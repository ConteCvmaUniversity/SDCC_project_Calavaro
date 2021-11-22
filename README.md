# SDCC_project_Calavaro - Marco Calavaro (matricola 0295233)

## Descrizione repository
Il repository contiene le varie parti richieste nella consegna del progetto di SDCC:
- La cartella `Application` contiene il codice sviluppato
- Il pdf riportato è la relazione prodotta
- Di seguito vi è la descrizione di howto del programma
	- Installazione
	- Configurazione
	- Esecuzione
### Configurazione iniziale
La configurazione iniziale che è riportata nel progetto prevede che allo start vengano istanziati:
- 3 peer      node
- 1 sequencer node
- 1 register  node

## HOWTO

### Installazione
Non vi sono particolari accorgimenti da prendere per poter eseguire l'applicazione, è necessario avere `docker-compose` installato.
Per poter eseguire il programma è necessario lanciare gli script dalla cartella `root` che in questo repository è `Application`

### Configurazione
Per poter eseguire la configurazione del progetto bisogna modificare i seguenti file:
- `.env` per modificare i parametri del docker-compose.yml
- `Application/utility/static_config.go` per modificare i parametri dell'applicazione
Salvate le opportune modifiche si potrà lanciare l'architettura tramite il comando:
```sh
./start.sh -b
``` 

#### Nota
Si è scelto di forzare l'utente a mantenere congruenti questi due file, questo perchè si è voluto adottare una soluzione che separasse l'applicativo Go da Docker, maggiori dettagli sono riportati nella relazione

### Esecuzione 
Per lanciare l'intera architettura basterà eseguire il comando
```sh
./start.sh 
``` 
Per potersi connettere ad un peer specifico (identificato da un numero PEERNUM) basterà eseguire il comando 
```sh
./connect_peer.sh PEERNUM
```


### Lancio dei test
Per poter eseguire i test bisogna modificare il valore di `Launch_Test` a `true` nel file `Application/utility/static_config.go` e lanciare il progetto ribuildando.
```sh
./start.sh -b
``` 