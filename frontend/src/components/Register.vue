<template>
  <div >    
    <p>
      Register-Tab
    </p>
    <div>  
        Es wurde kein gülltiges UserCookie gefunden.
        <br>
        Bitte registrieren
        <hr>
        Wer sind Sie? <input type="text" placeholder="Name" v-model="sendMessage.Name"/>
        <hr>
        Debug -> Name: {{this.sendMessage.Name}} UUID: {{this.sendMessage.Uuid}} 
        <button class="btn btn-secondary" v-on:click="RegisterIdent()">Register</button>
        <button class="btn btn-secondary" v-on:click="testMethod('bla')">TestMsg</button>
    </div>
  </div>
</template>

<script>
import { uuid } from 'vue-uuid'; 
import axios from 'axios';

/*
axios.interceptors.response.use(
  function(response) { return response;}, 
  function(error) {
    // handle error
    if (error.response) {
        //alert(error.response.data.message);
        console.log("axios interceptor: " + error.response.data);
        this.$parent.showAlert(error.response.data)
        
    }
  });

*/
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Register', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {          
          
          sendMessage: {
            MsgType: "RegisterRequest",            
            Name: "TestName",
            Uuid: "123456"
            }
          }
    },
  
  created() {

  },
  methods: {

    testMethod(m,str ) {
      console.log("here:" + m)
      this.$parent.showAlert(str)
    },
    async RegisterIdent() {
          console.log("RegisterIdent")
          this.sendMessage.Uuid = uuid.v1()
             axios.post(apiURL+ '/TestApi', this.sendMessage)
                 .then((res) => {
                     //Perform Success Action
                     console.log("Resut: "+ res)
                 })
                .catch((error) => {
                     console.log("Error:")
                     console.log("Error:"+ error.response.data.err)
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
