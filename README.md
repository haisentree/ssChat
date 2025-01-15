# ssChat
简单安全的Web聊天，服务端golang，通讯websocket，加密RSA

# 功能演示

0000000000（十个零）表示广播频道，里面的消息是不加密的。用户向广播发送消息，其他用户都可以收到，然后点击对应的消息可以添加客户端，被添加的一方也会自动添加。然后两者就可以实现RAS加密通讯了，需要切换到对应的通讯客户端才能够接收消息，因为没有设计消息存储功能。

![](https://raw.githubusercontent.com/haisentree/imageBed/main/image2024/QQ2025115-201316.gif)





