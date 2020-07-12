import datetime


def pretty_timedelta(module, delta):
    assert isinstance(delta, datetime.timedelta)
    s = []
    if delta.days:
        s.append(str(delta.days) + "d")

    seconds = delta.seconds
    if seconds > 3600:
        hr, seconds = divmod(seconds, 3600)
        s.append(str(hr) + "h")
    if seconds > 60:
        minutes, seconds = divmod(seconds, 60)
        s.append(str(minutes) + "m")
    if seconds:
        s.append(str(seconds) + "s")
    return "".join(s)
