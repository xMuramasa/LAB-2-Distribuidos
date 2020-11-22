# LAB-2-Distribuidos
2020-2
______________

## Integrantes:
    - Campus San Joaquin
    - 201773617-8   Zhuo Chang
    - 201773557-0   Martín Salinas
______________


## Maquinas 
Máquina 1   receiver --------> Tiene grpc & rabbit
hostname:   dist29
contraseña: BcSz2fUS

Máquina 2   camiones  --------> Tiene grpc
hostname:   dist31
contraseña: jzCsSjfR

Máquina 3   sender  --------> Tiene grpc
hostname:   dist30
contraseña: CtXTq9qq

Máquina 4   finanzas  --------> Tiene rabbit
hostname:   dist32
contraseña: k5PfFYfP

El usuario de las máquinas es: root
______________


## Instrucciones de uso:

#### Se debe entrar a la carpeta LAB-1-DS
#### Se debe compilar el archivo helloworld.proto dentro de la carpeta A_Dependencies usando el comando:
    protoc -I="." --go_out=$GOROOT/src/helloworld --go-grpc_out=$GOROOT/src/helloworld helloworld.proto

#### Maquina 1:

#### Maquina 2:

#### Maquina 3:

#### Maquina 4:
______________


## Consideraciones:
    - Se asume que hay un directorio en GOPATH llamado lab2
    - Se asume que el usuario sabe lo que tiene que hacer y no comete errores en sus inputs
    - Se asume que existe una carpeta llamada helloworld dentro de $GOROOT/src
    - Se asume que las variables de entorno estan correctamente actualizadas ($GOROOT, $GOPATH, $GOBIN) en el archivo, .bashrc ubicado en ~/ con lo siguiente
            export GOROOT=/usr/local/go
            export GOPATH=$HOME/go
            export GOBIN=$GOPATH/bin
            export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
