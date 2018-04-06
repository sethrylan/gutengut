const fetch = require('node-fetch');

'use strict'

exports.handler = (event, context, callback) => {
    var start = 0, 
        limit = 5;

    if (event.queryStringParameters) { 
      start = parseInt(event.queryStringParameters.start) || start;
      limit = parseInt(event.queryStringParameters.limit) || limit;
    }

  fetch('http://sethrylan.org/adventures.txt')
    .then(response => response.text())
    // .then(rawtext => rawtext.split('\n\n'))
    // .then(split => console.log(split.slice(1,2)))
    .then(rawtext => rawtext.split('\n\n').slice(start,start+limit).join('\n\n'))
    .then(text => callback(null, {
      statusCode: 200,
      headers: {"content-type": "text/plain"},
      body: text
    })
  );
};
