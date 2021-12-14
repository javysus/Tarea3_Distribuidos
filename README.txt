-Sistemas Distribuidos-
Laboratorio 3

Grupo 9:  Tomas Barber Villalobos            Rol: 201773591-0
          Pedro Mérida Álvarez               Rol: 201773583-K
          Javiera Villarroel Toloza          Rol: 201773580-5

Consideraciones:
    - El comando para Leia Organa "GetNumberRebelds" se dejo como tal, a pesar de creer que se produjo un typo en el enunciado
      del laboratorio (formalmente deberia ser "GetNumberRebels", sin la letra 'd')

Consistencia:
- El nodo dominante es el servidor 1, quien mantiene un temporizador para realizar el Merge cada 2 minutos.
- Los datos principales son los del nodo dominante (servidor 1), por ejemplo, si hay una misma ciudad en el servidor 1 y servidor 2 se mantiene el dato que está en el servidor 1.
- De la misma forma, si ocurren Updates de la misma ciudad en un servidor distinto del primero, se mantienen los datos del servidor 1.
- El orden a considerar para los comandos es: servidor 1, servidor 2 y servidor 3, por lo que si se agrega una misma ciudad en el servidor 2 y 3, se mantiene la del servidor 2. 
Es decir, la prioridad de los datos es en el orden servidor 1 - servidor 2 - servidor 3. Para lograrlo, se utiliza un arreglo con las ciudades agregadas, comenzando desde el servidor 1.

Ejecución:

- Para la correcta ejecucion del programa se deben ejecutar los archivos en el siguiente orden:
servidor_fulcrum3 - servidor_fulcrum2 - servidor_fulcrum1 - broker_mos_eisley - ahsoka - thrawn - leia
- Abrir los puertos 5051, 5052, 5053 y 5054 para cada máquina virtual. Nosotros utilizamos el siguiente comando:
sudo firewall-cmd --zone=public --add-port=50052/tcp --permanent (para cada puerto)

Como ejecutar el servidor fulcrum 3:
- Ubicarse en la maquina dist37
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_server3", luego ejecutar el comando "make"

Cómo ejecutar el servidor fulcrum 2:
- Ubicarse en la maquina dist38
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_server2", luego ejecutar el comando "make"

Cómo ejecutar el servidor fulcrum 1:
- Ubicarse en la maquina dist39
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_server", luego ejecutar el comando "make"

Cómo ejecutar el broker:
- Ubicarse en la maquina dist40
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_broker, luego ejecutar el comando "make"

Cómo ejecutar el Informante Ahsoka:
- Ubicarse en la máquina dist37
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_client1", luego ejecutar el comando "make"

Cómo ejecutar el Informante Thrawn:
- Ubicarse en la máquina dist38
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_client2", luego ejecutar el comando "make"

Cómo ejecutar Leia:
- Ubicarse en la máquina dist39
- Ubicarse en la carpeta "Tarea3_Distribuidos" (cd Tarea3_Distribuidos)
- Ubicarse en la carpeta "rebelion_leia", luego ejecutar el comando "make"
