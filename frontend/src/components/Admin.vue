<template>
  <div >    
    <p>
      Admin-Tab
    </p>
    <div>  
        <div>
            <button class="btn btn-secondary" v-on:click="GetUsers()">GetUsers</button>
            <br> Debug: Users: {{this.Users}}


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


export default {
  name: 'Admin', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {
          Users: [{}]
      }
  },
  created() {

  },
  methods: {
    async GetUsers() {
          console.log("GetUsers from DB")
      try {
        
        let response = await axios.get(apiURL  + "/AdminGetUsers");
        this.Users = response.data;
        console.log(response.status)
        if (response.status != 200 ) {        
            console.log("Users not object")
            this.$parent.showAlert("Server returned an Error:" +response.status+ " - "+ response ); 
        }
      } catch (error) {
        console.log(error);
        
        this.$parent.showAlert("Server returned an Error:\n" + error); 
      }
    },
    ChangeUserEnabled() {
        //ToDo: api call here
        this.GetUsers()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
