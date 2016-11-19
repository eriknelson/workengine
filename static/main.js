'use strict'

$(function() {
  main()
});

function main(socket) {
  socket = io();

  socket.on('connect', function() {
    console.log('Connected')
  })

  socket.on('firehose', function(msg) {
    writeToViewer(msg);
  })
}

function genUrl(path) {
  return 'http://localhost:3000' + path;
}

function writeToViewer(msg) {
  $('.log-viewer').append('<div>' + msg + '</div>');
  $('.log-viewer').animate({
    scrollTop: $('.log-viewer').prop("scrollHeight")
  }, 0);
}
