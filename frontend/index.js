const app = new Vue({
    el: '#app',
    vuetify: new Vuetify(),
    data: {
        header:[
          "ID",
          "Type",
          "ToDo"
        ],
        todos: [],
        AddToDo: {}
      },
      methods: {
        getTodos() {
          axios.get("http://localhost:8080/todos", {}).then((response) => {
            this.todos = response.data;
          }).catch((error) => {
            console.log(error);
          })
        },
        createToDo(){
          const todoData = Object.assign ({}, this.AddToDo)
          axios.post("http://localhost:8080/todos", todoData).then((response) => {
            alert("登録しました")
            this.getTodos()
            this.AddToDo = {}
          }).catch((error) => {
            console.log(error);
          })
        },
        deleteToDo(id) {
          axios.delete("http://localhost:8080/todos", {
            params: {
              id: id,
            }
          }).then((response) => {
            alert("削除しました")
            this.getTodos()
          }).catch((error) => {
            console.log(error)
          })
        }
      },
      mounted(){
        this.getTodos();
      }
  })