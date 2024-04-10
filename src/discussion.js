//number id ——> uid
const express = require('express');
const router = express.Router();
const moment = require('moment-timezone');

let posts = [];
let replies = [];

// 创建帖子
router.post('/posts', (req, res) => {
    const { userId, content } = req.body;
    if (!userId || !content) {
        return res.status(400).send('UserId and content are required');
    }
    const newPost = {
        id: posts.length + 1,
        userId,
        content,
        timestamp: moment().tz('America/New_York').format(),
        //timestamp: new Date(),
    };
    posts.push(newPost);
    res.status(201).send(newPost);
});

// 对帖子进行回复
router.post('/posts/:postId/replies', (req, res) => {
    const { userId, content } = req.body;
    const { postId } = req.params;

    if (!userId || !content) {
        return res.status(400).send('UserId and content are required');
    }

    const postExists = posts.some(post => post.id === parseInt(postId));
    if (!postExists) {
        return res.status(404).send('Post not found');
    }

    const newReply = {
        id: replies.length + 1,
        postId: parseInt(postId),
        userId,
        content,
        timestamp: moment().tz('America/New_York').format(),
        //timestamp: new Date(),
    };
    replies.push(newReply);
    res.status(201).send(newReply);
});

// 获取所有帖子
router.get('/posts', (req, res) => {
    res.send(posts);
});

// 获取单个帖子
router.get('/posts/:postId', (req, res) => {
    const { postId } = req.params;
    const post = posts.find(post => post.id === parseInt(postId));
    if (!post) {
        return res.status(404).send('Post not found');
    }
    res.send(post);
});

// 获取单个帖子及其所有回复
router.get('/posts/:postId/replies', (req, res) => {
    const { postId } = req.params;
    const post = posts.find(post => post.id === parseInt(postId));
    if (!post) {
        return res.status(404).send('Post not found');
    }

    const postReplies = replies.filter(reply => reply.postId === parseInt(postId));

    res.send({
        post,
        replies: postReplies
    });
});


// 获取特定帖子的所有回复
// router.get('/posts/:postId/replies', (req, res) => {
//     const { postId } = req.params;
//     const postReplies = replies.filter(reply => reply.postId === parseInt(postId));
//     res.send(postReplies);
// });

module.exports = router;