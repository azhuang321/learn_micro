<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: helloworld.proto

namespace Greeter;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>Greeter.HelloRequest</code>
 */
class HelloRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string name = 1;</code>
     */
    protected $name = '';
    /**
     * Generated from protobuf field <code>repeated string hobby = 2;</code>
     */
    private $hobby;
    /**
     * Generated from protobuf field <code>.Greeter.Sex sex = 3;</code>
     */
    protected $sex = 0;
    /**
     * Generated from protobuf field <code>map<string, string> Mp = 4;</code>
     */
    private $Mp;

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $name
     *     @type string[]|\Google\Protobuf\Internal\RepeatedField $hobby
     *     @type int $sex
     *     @type array|\Google\Protobuf\Internal\MapField $Mp
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Helloworld::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string name = 1;</code>
     * @return string
     */
    public function getName()
    {
        return $this->name;
    }

    /**
     * Generated from protobuf field <code>string name = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setName($var)
    {
        GPBUtil::checkString($var, True);
        $this->name = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>repeated string hobby = 2;</code>
     * @return \Google\Protobuf\Internal\RepeatedField
     */
    public function getHobby()
    {
        return $this->hobby;
    }

    /**
     * Generated from protobuf field <code>repeated string hobby = 2;</code>
     * @param string[]|\Google\Protobuf\Internal\RepeatedField $var
     * @return $this
     */
    public function setHobby($var)
    {
        $arr = GPBUtil::checkRepeatedField($var, \Google\Protobuf\Internal\GPBType::STRING);
        $this->hobby = $arr;

        return $this;
    }

    /**
     * Generated from protobuf field <code>.Greeter.Sex sex = 3;</code>
     * @return int
     */
    public function getSex()
    {
        return $this->sex;
    }

    /**
     * Generated from protobuf field <code>.Greeter.Sex sex = 3;</code>
     * @param int $var
     * @return $this
     */
    public function setSex($var)
    {
        GPBUtil::checkEnum($var, \Greeter\Sex::class);
        $this->sex = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>map<string, string> Mp = 4;</code>
     * @return \Google\Protobuf\Internal\MapField
     */
    public function getMp()
    {
        return $this->Mp;
    }

    /**
     * Generated from protobuf field <code>map<string, string> Mp = 4;</code>
     * @param array|\Google\Protobuf\Internal\MapField $var
     * @return $this
     */
    public function setMp($var)
    {
        $arr = GPBUtil::checkMapField($var, \Google\Protobuf\Internal\GPBType::STRING, \Google\Protobuf\Internal\GPBType::STRING);
        $this->Mp = $arr;

        return $this;
    }

}
