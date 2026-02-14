package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

// UpdateMatterByFile expects format 1234-2020
func (s *SyncApp) UpdateMatterByFile(q string) error {
	ctx := context.Background()
	file := q

	if !strings.Contains(q, " ") {
		file = fmt.Sprintf("Int %s", q)
	}
	fileType, _, _ := strings.Cut(file, " ")
	switch fileType {
	case "Int":
		fileType = "Introduction"
	case "Res":
		fileType = "Resolution"
	case "LU":
		fileType = "Land Use Application"
	default:
		return fmt.Errorf("unknown legislation type for %q", file)
	}
	filter := legistar.AndFilters(
		legistar.MatterTypeFilter(fileType),
		legistar.MatterFileFilter(file),
	)

	matters, err := s.legistar.Matters(ctx, filter)
	if err != nil {
		return err
	}
	if len(matters) != 1 {
		return fmt.Errorf("expected 1 response got %d for %q", len(matters), q)
	}
	switch fileType {
	case "Introduction":
		return s.UpdateLegislation(ctx, matters[0].ID)
	case "Resolution":
		return s.UpdateResolution(ctx, matters[0].ID)
	case "Land Use Application":
		return s.UpdateLandUse(ctx, matters[0].ID)
	default:
		return fmt.Errorf("unknown legislation type for %q", file)
	}
}

func (s *SyncApp) updateMatter(ctx context.Context, fn string, l db.Legislation) error {
	sponsors, err := s.legistar.MatterSponsors(ctx, l.ID)
	if err != nil {
		return err
	}
	l.Sponsors = []db.PersonReference{}
	for _, p := range sponsors {
		if p.MatterVersion != l.Version {
			continue
		}
		l.Sponsors = append(l.Sponsors, s.personLookup[p.NameID].Reference())
	}

	history, err := s.legistar.MatterHistories(ctx, l.ID)
	if err != nil {
		return err
	}
	l.History = nil
	for _, mh := range history {
		hh := db.NewHistory(mh)
		// For some older entries spuriously have PassedFlagName=Fail where the API call to get votes errors
		// 27 == "Introduced by Council"
		// 53 == "Sent to Mayor by Council"
		// 5014 == "Hearing Held by Mayor"
		// 5023 == "Recved from Mayor by Council"
		if hh.PassedFlagName != "" && hh.ActionID != 53 && hh.ActionID != 5014 && hh.ActionID != 5023 && hh.ActionID != 27 {
			votes, err := s.legistar.EventVotes(ctx, hh.ID)
			if err != nil {
				switch {
				case hh.PassedFlagName == "Failed" || hh.ID == 283777:
					log.Printf("warning getting votes for eventID %v - PassedFlagName:Failed is a known bug on older records", hh.ID)
				case hh.ID == 433502:
					// https://bsky.app/profile/jehiah.cz/post/3meqwcebuc22v
					log.Printf("warning getting votes for eventID %v - known bug ", hh.ID)
				default:
					return err
				}
			}
			hh.Votes = db.NewVotes(votes)
		}
		l.History = append(l.History, hh)
	}

	attachments, err := s.legistar.MatterAttachments(ctx, l.ID)
	if err != nil {
		return err
	}
	l.Attachments = nil
	for _, a := range attachments {
		l.Attachments = append(l.Attachments, db.NewAttachment(a))
	}

	versions, err := s.legistar.MatterTextVersions(ctx, l.ID)
	if err != nil {
		return err
	}
	l.TextID = versions.LatestTextID()
	txt, err := s.legistar.MatterText(ctx, l.ID, l.TextID)
	if err != nil {
		return err
	}
	l.Text = txt.SimplifiedText()
	l.RTF = txt.SimplifiedRTF()

	if len(l.RTF) > 51200 { // 50k
		l.RTF = ""
	}
	return s.writeFile(fn, l)
}
