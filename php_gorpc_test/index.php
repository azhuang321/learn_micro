<?php
include "vendor/autoload.php";


// 用于连接 服务端
$client = new \GreeterClient\GreeterClient('127.0.0.1:8000', [
    'credentials' => Grpc\ChannelCredentials::createInsecure(),
]);


// 实例化 GetUserRequest 请求类
$request = new \Greeter\HelloRequest();
$request->setName("test");
$request->setHobby(["swimming","running"]);
$request->setSex(\Greeter\Sex::FEMALE);
$request->setMp([
    "name" => "test1",
    "sex" => "test2",
]);

// 调用远程服务
$get = $client->SayHello($request,[
    "username" => ["test1"],
    "password" => ["test2"]
])->wait();

// $reply  是 SayHello 返回对象
// $status 是 记录 grpc 错误信息 对象
list($reply, $status) = $get;

var_dump($reply->getMessage());
var_dump($reply->getHobby()->offsetGet(1));
var_dump($reply->getSex());
var_dump($reply->getMp()->offsetGet("name"));
//var_dump($status);