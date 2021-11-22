# SDCC_Project - Marco Calavaro (matricola 0295233)

## Esecuzione del progetto

### Avvio dell'applicazione
Per poter avviare l'applicazione basta eseguire nella `root` del progetto i comandi di `docker-compose`, per semplificarne l'utilizzo sono stati creati i bash script `start.sh` e `stop.sh` eseguibili nei seguenti modi:
- Per effetuare l'up:
```sh
./start.sh 
``` 
- Ribuildando il progetto (up --build):
```sh
./start.sh -b
``` 
- Per effetuare il down:
```sh
./stop.sh 
``` 
### Collegarsi a un peer 
Per poter accedere ad un peer bisogna posizionarsi nella cartella `root` del progetto ed eseguire il comando
```sh
./connect_peer.sh PEERNUM
```

### Lancio dei test
Per poter eseguire i test bisogna modificare il valore di `Launch_Test` a `true` nel file `SDCC_Project/utility/static_config.go` e lanciare il progetto ribuildando.
```sh
./start.sh -b
``` 
