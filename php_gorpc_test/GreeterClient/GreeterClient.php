<?php
namespace GreeterClient;

class GreeterClient extends \Grpc\BaseStub
{
    public function __construct($hostname, $opts, $channel = null)
    {
        parent::__construct($hostname, $opts, $channel);
    }

    public function SayHello(\Greeter\HelloRequest $argument,$metadata = [],$options = []) {

        return $this->_simpleRequest(
            '/Greeter/SayHello',//远端服务
            $argument,
            ['\Greeter\HelloReply', 'decode'],//返回解码
            $metadata,
            $options
        );
    }

}