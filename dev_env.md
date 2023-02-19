# Development environment for taskPlantDep

<!-- vim: set conceallevel=0 :-->

Project is developed using neovim.
So you will find some vim annotations from time to time.

## live testing

In a first terminal launch :

```bash
java -jar tplantuml.jar
```

In a second terminal use following command:

```bash
ls  *.go model/*.go cleaning/*.go | entr go run taskPlant.go
```
