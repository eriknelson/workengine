'use strict'

$(function() {
  main()
});

function main() {
  console.log('Hello wworld!');
  axios.post(genUrl('/run'), {
    id: uuid.v4()
  }).then(res => {
    console.log('got response!');
    console.log(res.data);
  });
}

function genUrl(path) {
  return 'http://localhost:3000' + path;
}
