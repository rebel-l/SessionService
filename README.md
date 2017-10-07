[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Dependency Status](https://www.versioneye.com/user/projects/591d7054d83ae5005cde5b7d/badge.svg?style=flat-square)](https://www.versioneye.com/user/projects/591d7054d83ae5005cde5b7d) [![Build Status](https://travis-ci.org/rebel-l/SessionService.svg?branch=develop)](https://travis-ci.org/rebel-l/SessionService)

# Session Service
This service delivers several endpoints to create, load, change and delete sessions.

# Requirements
## <a name="reqman"></a>Mandatory
The only requirement so far is [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/). 

## <a name="reqopt"></a>Optional
Optionally you are able to run the whole environment in a virtual machine. Therefor you need:
* [Vagrant](https://www.vagrantup.com/) with plugins:
	* vagrant-hostmanager
	* vagrant-vbguest
* [Oracle Virtual Box](https://www.virtualbox.org/)
* [PHP](http://www.php.net/)
* [Composer](https://getcomposer.org/)  

# Development
To get the development environment run you only need to follow the instructions under [Docker Environment](#dockerenv).
If you decided to run it on a virtual machine, then please do the steps in [Vagrant Environment](#vagrantenv). 

## <a name="dockerenv"></a>Docker Environment
Before you can start ensure that you have install all [requirements](#reqman). All 
commands will be executed in the projects _root_ folder.

### Docker Compose Way
The easiest way is to use docker compose. You are able to launch every container with only one command:
```bash
docker-compose up -d
```

To work with go and the addional tools, you can access the container with:
```bash
docker-compose exec sessionservice ./docker-entrypoint
```

On the command line of the docker container you can now find all the projects data at:
```bash
cdproj
```

### Docker Way
Now you can build the docker image by excuting:
```bash
docker build -t sessionservice .
```

Afterwards you can run the docker container by:
```bash
docker run -it -p 4000:4000 --name sessionservice -v /vagrant/:/workspace/src/github.com/rebel-l/sessionservice sessionservice
```

On the command line of the docker container you can now find all the projects data at:
```bash
cdproj
```

The Golang compiler should be able to execute from everywhere. You can check that by:
```bash
go version
``` 

To detach from the docker container you need just to press the keys _Ctrl + p Ctrl + q_.

To launch the Redis container execute the following:
```bash
docker run -it -p 6379:6379 --name redis -d redis redis-server --appendonly yes
```

## <a name="vagrantenv"></a>Vagrant Environment 
If you would like to have clean sandbox for everything a vagrant machine is maybe your choice. 
Therefor you need to ensure to have the [optional requirements](#regopt) ready. Ensure that you
run all commands in the projects _root_ folder.

To install all necessary packages just run _composer_:
```bash
composer install
```

Afterwards you are able to start the virtual machine:
```bash
vagrant up
```

Now you can connect with your favourite _ssh_ tool to the virtual machine:
```bash
vagrant ssh
```
_You can use the dns name 'session.dev' for connection_.

You can switch to the project folder by:
```bash
cd /vagrant # or the alias 'cdproj'
```

On your virtual machine you can run your docker like described in the [Docker Environment](#dockerenv).

# Quality Assurance
For quality assurance [Travis CI](https://travis-ci.org) is connected with this repository. But before committing or 
pushing anything to this repository you can quickly check everything by executing the build script:
```bash
./scripts/build.sh
```

# API Documentation
You find actual version of the API documentation in the _/docs/indext.html_ file of this repository or as an online 
version on [Swagger hub](https://app.swaggerhub.com/apis/rebel-l/SessionService).

The swagger file itself can be also found as YAML version in the _/docs_ subfolder of this repository.

If the service is released you can also open http://servicedomain/docs/, e.g. on the vagrant development environment it is
http://session.dev:4000/docs/. 