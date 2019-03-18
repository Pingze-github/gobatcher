import tornado.httpserver
import tornado.ioloop
import tornado.web
import tornado.gen
import time

class MainHandler(tornado.web.RequestHandler):
    @tornado.web.asynchronous
    @tornado.gen.engine
    def get(self):
        yield tornado.gen.sleep(0.01)
        self.write("")
        self.finish()

if __name__ == "__main__":
    app = tornado.web.Application(handlers=[(r"/", MainHandler)])
    http_server = tornado.httpserver.HTTPServer(app)
    http_server.listen(9003)
    tornado.ioloop.IOLoop.instance().start()