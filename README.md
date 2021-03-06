# `subtitle`
`subtitle`是用于获取字幕的工具目前支持射手，迅雷两个平台获取字幕

[![Build Status](https://app.travis-ci.com/YouCD/subtitle.svg?branch=main)](https://app.travis-ci.com/YouCD/subtitle)
# 安装
* 自行编译
```shell
git clone install github.com/YouCD/subtitle
cd subtitle
make build&&sudo mv bin/subtitle/subtitle /usr/local/bin
```
* 快速安装
```shell
go install github.com/YouCD/subtitle@latest

```
# 食用方法
`subtitle`默认从射手网往下载字幕

* 自动获取当前目下视频文件，并下载字幕，字幕源为迅雷
```shell
subtitle --source xunlei 
```
* 自动获取当前目下视频文件，并下载字幕，字幕源为射手网
```shell
subtitle 
```
* 指定目录
```shell
# subtitle  -d Forever.US.S01 
The subtitle file Forever.US.S01/Forever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.ass is saved.
The subtitle file Forever.US.S01/Forever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.01.srt is saved.
The subtitle file Forever.US.S01/Forever.2014.S01E01.Pilot.1080p.WEB-DL.DD5.1.H.264-ECI.02.srt is saved.

```
* 获取帮助
```shell
subtitle --help            
subtitle 是用于获取字幕的工具

Usage:
  subtitle [flags]
  subtitle [command]

Examples:
    指定电影文件下载
        subtitle -f SomeVideo.mkv
    指定电影文件夹
        subtitle -d SomeVideDir
    指定源 支持 Shooter Xunlei  OpenSubtitles(可简写为:osb)   默认为射手网 Shooter
        subtitle --source xunlei -d SomeVideDir
        subtitle --source  osb   -d SomeVideDir

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version info of subtitle

Flags:
  -d, --dir string      指定电影存放目录
  -f, --file string     指定电影文件,如果不指定则使用当前目录
  -h, --help            help for subtitle
      --source string   指定字幕源,如果不指定则默认为 shooter (default "shooter")

Use "subtitle [command] --help" for more information about a command.

```