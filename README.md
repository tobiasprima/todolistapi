# Golang Restful API with MongoDB Database
This project is hosted on https://todolistapi-kcjj.onrender.com

This is a backend for TodoList Web Apps on
https://github.com/tobiasprima/todolist_react_web

## Quick Start
1. Clone

```bash
git clone https://github.com/tobiasprima/todolist_react_web
```

2. Install Dependency

```bash
go mod tidy
```

3. To test with Mongodb locally
-create new `.env` file and add:

```bash
DATABASE_URI = {YOUR_DATABASE_URI}
```

-run on localhost

```bash
go run main.go
```

4. Test Endpoints
on postman

### get todos (GET)
```bash
http://localhost:8080/todos
```

### create new todo (POST)
```bash
http://localhost:8080/todo
```
with body
```bash
{
  "title" : {INPUT_TITLE}
}
```
Status will default to false (Uncompleted)

### update todo status (PATCH)
```bash
http://localhost:8080/todo/{ID}   // copy ID from get/todos method
```


### reorder todo based on status (PATCH)
This will create an order property for todos and it will automatically sorted in the get method
```bash
http://localhost:8080/todos/reorder
```

### reset todo order (POST)
```bash
http://localhost:8080/todos/reset
