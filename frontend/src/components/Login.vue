<template>
  <div >    
    <p>
      Login-Tab
    </p>
    <div>  
        Es wurde kein gülltiges UserCookie gefunden.
        <br>
        Bitte registrieren
        <hr>
        Wer sind Sie? <input type="text" placeholder="Name" v-model="sendMessage.Name"/>
        <hr>
        Debug -> Name: {{this.sendMessage.Name}} UUID: {{this.sendMessage.Uuid}} 
        <button class="btn btn-secondary" v-on:click="LoginRequest()">Login</button>
        
    </div>
  </div>
</template>

<script>

import axios from 'axios';

const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Login', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {          
          
          sendMessage: {
            MsgType: "LoginRequest",            
            Name: ""            
            }
          }
    },
  
  created() {

  },
  methods: {
    
    async LoginRequest() {
          console.log("LoginRequest")
            
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
                 });
    }
  }        
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
