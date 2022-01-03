import tornado.options
import tornado.httpserver
import tornado.httpclient
import tornado.ioloop
import tornado.web
import os.path
import logging
import os
import urllib.parse
import json

import uimodules
import uimethods

class BaseHandler(tornado.web.RequestHandler):
    http = tornado.httpclient.AsyncHTTPClient()

    async def legislator_req(self, path):
        token = os.environ.get("NYC_LEGISLATOR_TOKEN")
        BASE="https://webapi.legistar.com/v1/nyc"
        url = BASE + "/" + path
        if self.request.arguments:
            url += "?" + urllib.parse.urlencode(self.request.arguments, doseq=True) #.replace("%24", "$")
        if "?" in url:
            url += "&token=" + urllib.parse.quote(token)
        else:
            url += "?token=" + urllib.parse.quote(token)
        logging.info('url %r', url)
        resp =  await self.http.fetch(url, method="GET")
        return resp


class Index(BaseHandler):
    async def get(self, name=None):
        if not name:
            self.render("index.html", name=None)
        else:
            resp = await self.legislator_req(name)
            logging.info('name %r', name)
            # logging.info('response %r', resp.body)
            if self.get_argument("format",None) == "json":
                self.write(resp.body)
                return
            self.render("index.html", name=name, data=json.loads(resp.body))

class Application(tornado.web.Application):
    def __init__(self, testing=False):
        app_settings = {
            'debug': True,
            "template_path": os.path.join(os.path.dirname(__file__), "templates"),
            "static_path": os.path.join(os.path.dirname(__file__), "static"),
            "ui_modules": uimodules,
            "ui_methods": uimethods,
#            "autoescape": True,
            }
        handlers = [
            (r"^/(.*)$", Index),
        ]
        tornado.web.Application.__init__(self, handlers, **app_settings)


if __name__ == "__main__":
    tornado.options.define("port", default=7001, help="Listen on port", type=int)
    tornado.options.parse_command_line()
    http_server = tornado.httpserver.HTTPServer(Application(), xheaders=True)
    http_server.listen(tornado.options.options.port, address="0.0.0.0")
    logging.info("listening on 0.0.0.0:%s", tornado.options.options.port)
    tornado.ioloop.IOLoop.instance().start()
