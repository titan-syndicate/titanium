const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => {
  res.send('Hello from test Node.js app!');
});

app.listen(port, () => {
  console.log(`Test app listening at http://localhost:${port}`);
});