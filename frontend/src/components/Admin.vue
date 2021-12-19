<template>
  <div >    
    <p>
      Admin-Tab
    </p>
    <div>  
      <div>
        Neuer Username: <input type="text" placeholder="Name" v-model="sendMessage.Name"/>
     <!--   Neue Uuid: <input type="text" placeholder="Uuid" v-model="sendMessage.Uuid"/> -->
        <button class="btn btn-secondary" v-on:click="AddUser()">Create new user</button>
      </div>
        <br>
        <div>
            <button class="btn btn-secondary" v-on:click="GetUsers()">GetUsers</button>
          <!--  <br> Debug: Users: {{this.Users}} --> 


            <table class="table table-dark">
              <thead>
                <tr>            
                  <th scope="col">Username</th>
                  <th scope="col">Enabled</th>
                  <th scope="col">Change</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in this.Users.User" :key="item.Name">   
                  <td>{{ item.Name }}</td>
                  <td>{{ item.Enabled }}</td>      
                  <td><button class="btn btn-secondary" v-on:click="ChangeUserEnabled()">Change</button></td>      
                </tr>
              </tbody>            
            </table>
        </div>
    </div>
  </div>
</template>

<script>

import axios from 'axios';
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"
//import { uuid } from 'vue-uuid'; 
// import ShortUniqueId from 'short-unique-id';


export default {
  name: 'Admin', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {
          Users: [{}],
          sendMessage: {
            MsgType: "AddUserRequest",            
            Name: "",
            Uuid: ""
            }
      }
  },
  created() {
    //const uid = new ShortUniqueId();

    //this.sendMessage.Uuid = uuid.v4()
    //this.sendMessage.Uuid = uid()
    this.GetUsers()
  },
  methods: {
    async GetUsers() {
        console.log("GetUsers")
      this.sendMessage.MsgType = "GetUsers"           

             axios.post(apiURL+ '/TestApi', this.sendMessage)
                 .then((res) => {
                     //Perform Success Action
                     console.log("Resut: "+ res.data.Result)
                     this.Users = res.data;
                 })
                .catch((error) => {                     
                     console.log("Error:"+ error.response.data.Result)
                     this.$parent.showAlert(error.response.data)
                })
                 .finally(() => {
                     //Perform action in always
                 });
        
  
    },
    ChangeUserEnabled() {
        //ToDo: api call here
        this.GetUsers()
    },
    async AddUser() {       
          console.log("AddUser")
          if (this.sendMessage.Name =="") {
            this.$parent.showAlert("Bitte einen Benutzername eingeben!")
          }else {
            this.sendMessage.MsgType = "AddUserRequest"           

             axios.post(apiURL+ '/TestApi', this.sendMessage)
                 .then((res) => {
                     //Perform Success Action
                     console.log("Resut: "+ res.data.Result)
                 })
                .catch((error) => {                     
                     console.log("Error:"+ error.response.data.Result)
                     this.$parent.showAlert(error.response.data)
                })
                 .finally(() => {
                     //Perform action in always
                    this.GetUsers()           
                 });
          }
          
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
