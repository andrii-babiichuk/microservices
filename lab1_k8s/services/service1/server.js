const http = require('http');

const requestListener = function (req, res) {
    if (req.url === '/api/service1') {
        res.writeHead(200);
        res.write('Hello from node server');
    } else {
        res.writeHead(404);
    }
    res.end()
}

const server = http.createServer(requestListener);
server.listen(8080);
