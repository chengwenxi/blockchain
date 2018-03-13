# gaia
A demonstration of the Cosmos-Hub with basic staking functionality.

+ [gaia简介](#gaia简介)
  + [gaia是什么](#gaia是什么)
  + [gaia架构](#gaia架构)
  + [gaia功能](#gaia功能)
    + [gaia client](#client)
    + [gaia node](#node)
    + [gaia rest-server](#rest-server)
  + [gaia代码结构](#gaia代码结构)
  + [gaia更新](#gaia更新)

+ [gaia私链搭建及配置详解](#gaia私链搭建及配置文件详解)
  + [私链搭建](#gaia私链搭建)
  + [配置详解](#配置文件详解)
  + [功能演示](#功能演示)


## gaia简介

## gaia是什么
  + gaia是cosmos一个project——具有基本的staking功能的Cosmos-Hub的示范。    https://github.com/cosmos/gaia

  + gaia也是cosmos-hub目前的测试网络使用名，官方发布过两个测试网络gaia1和gaia2。

#### gaia架构

#### gaia功能
  + 功能架构
  ![img](./source/gaia架构.png)

  + 运行时状态
  ![img](./source/运行时状态.png)

##### client

##### node

#### rest-server

#### gaia代码结构

https://github.com/cosmos/gaia

#### gaia更新

[cosmos roadmap](https://cosmos.network/roadmap)

gaia 依赖关系

| branch  | version | cosmos-sdk | tendermint |  
| ------- |:-------:| -----:     | -----:     |          
| master  | 0.5.0   | develop    | v0.15.0    |
| develop | 0.6.0   | tm-develop | develop    |


CHANGELOG
+ [tendermint](https://github.com/tendermint/tendermint/blob/master/CHANGELOG.md#0160-february-20th-2017)

+ [cosmos-sdk](https://github.com/cosmos/cosmos-sdk/blob/master/CHANGELOG.md)

+ [gaia](https://github.com/cosmos/gaia/blob/master/CHANGELOG.md)

## gaia私链搭建及配置文件详解

#### gaia私链搭建
  + [安装及私链搭建](Local-Test)


#### 配置文件详解
  + 文件结构
  + [priv_validator.json](config/priv_validator.json) 和 [genesis.json](config/genesis.json)
  + [config.toml](config/config.toml)

#### 功能演示
  + gaia delegate演示
