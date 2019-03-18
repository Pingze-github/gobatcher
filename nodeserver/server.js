const http = require('http');

const server = http.createServer(function (request, response) {
  return setTimeout(() => {
    response.write('')
    response.end()
  }, 10)
})

console.log('Server running ...')
server.listen(9002, '127.0.0.1')