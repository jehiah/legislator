import tornado.web
import tornado.escape

class RecursiveDump(tornado.web.UIModule):

    def render(self, data, key=None):
        if isinstance(data, (int, float)):
            return str(data)
        if isinstance(data, (str)):
            return tornado.escape.xhtml_escape(data)
        if isinstance(data, None.__class__):
            return 'None'

        # if isinstance(data, dict):
        #     for k in data:
        #         for whitelist in ["a"]:
        #             if whitelist in k:
        #                 data[k] = "REDACTED"

        return self.render_string("recursive-dump.html",
                                  data=data,
                                  key=key,
                                  timestamp_fields=set(['activated_ts', 'modified_ts', 'created_ts', 'deactivated_ts']),
                                  guid_fields=set(['guid', 'PersonGuid', 'MatterGuid', 'OfficeRecordGuid', 'EventGuid', 'IndexGuid', 'VoteGuid']),
                                  )
