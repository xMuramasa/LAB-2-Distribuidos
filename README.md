# LAB-2-Distribuidos
2020-2
______________

## Integrantes:
    - Campus San Joaquin
    - 201773617-8   Zhuo Chang
    - 201773557-0   Martín Salinas
______________


## Maquinas 
Máquina 1       (NameNode)   
hostname    :   dist29
contraseña  :   BcSz2fUS

Máquina 2       (DataNode1) 
hostname    :   dist31
contraseña  :   jzCsSjfR

Máquina 3       (DataNode2 && Clients) 
hostname    :   dist30
contraseña  :   CtXTq9qq

Máquina 4       (DataNode3) 
hostname    :   dist32
contraseña  :   k5PfFYfP

El usuario de las máquinas es: root
______________


## Instrucciones de uso:

#### Se debe entrar a la carpeta LAB-2-DISTRIBUIDOS
#### Se debe compilar el archivo SendReceive.proto dentro de la carpeta AA_Dependencies usando el comando:
    - make
### Antes de todo, en el Namenode y datanodes, se debe seleccionar el algoritmo a utilizar, ingresar 1 o 2 según se quiera!

#### Maquina 1 NAMENODE:
    - Asume que hay datanodes funcionando
    - Para inicializar el namenode, se debe ingresar al directorio NameNode y ejecutar el comando make
    - Para detener el servicio basta con ctrl+c
    - Para limpiar el registro: make clean

#### Maquina 2:
    - Asume que el NameNode está funcionando
    - Para inicializar el Datanode 1, se debe ingresar al directorio DataNodes/DN1 y ejecutar el comando make
    - Para detener el servicio basta con ctrl+c
    - Para limpiar los archivos almacenados dentro del directorio stored/, ejecutar make clean
#### Maquina 3:
  - DataNode 2:
    - Asume que el NameNode está funcionando
    - Para inicializar el Datanode 2, se debe ingresar al directorio DataNodes/DN2 y ejecutar el comando make
    - Para detener el servicio basta con ctrl+c
    - Para limpiar los archivos almacenados dentro del directorio stored/, ejecutar make clean
  
  - ClientDownloader:
    - Asume que el NameNode está funcionando
    - Ejecutar con make
    - Se debe ingresar el nombre exacto del libro a descargar con guiones y extención ej: CONAN_EL_VENGADOR-Robert_E._Howard.pdf
  
  - ClientUploader:
    - Asume que el NameNode está funcionando
    - Ejecutar con make
    - Se debe ingresar el nombre exacto del libro a subir con guiones y extención ej: CONAN_EL_VENGADOR-Robert_E._Howard.pdf
#### Maquina 4:
    - Asume que el NameNode está funcionando
    - Para inicializar el Datanode 3, se debe ingresar al directorio DataNodes/DN3 y ejecutar el comando make
    - Para detener el servicio basta con ctrl+c
    - Para limpiar los archivos almacenados dentro del directorio stored/, ejecutar make clean
______________


## Consideraciones:
    - Se asume que siempre hay al menos un datanode funcionando y el namenode está activo
    - Si el clientUploader falla la conneccion inicial, lo intentará denuevo hasta establecer una conneccion con otro dn
    - Se asume que hay un directorio en GOPATH llamado lab2
    - Se asume que el usuario sabe lo que tiene que hacer y no comete errores en sus inputs
    - Se asume que existe una carpeta llamada lab2 dentro de $GOROOT/src
    - Se asume que las variables de entorno estan correctamente actualizadas ($GOROOT, $GOPATH, $GOBIN) en el archivo, .bashrc ubicado en ~/ con lo siguiente
            export GOROOT=/usr/local/go
            export GOPATH=$HOME/go
            export GOBIN=$GOPATH/bin
            export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
