const express = require('express');
const discussionRoutes = require('./discussion');
const app = express();
const PORT = 3000;

app.use(express.json());

app.use('/discussion', discussionRoutes);

app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});