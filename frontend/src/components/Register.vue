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
        Wer sind Sie? <input type="text" placeholder="Name" v-model="Name"/>
        <hr>
        Debug -> Name: {{this.Name}} UUID: {{this.Uuid}} Status: {{this.Status}}
        <button class="btn btn-secondary" v-on:click="RegisterIdent()">Register</button>
    </div>
  </div>
</template>

<script>
import { uuid } from 'vue-uuid'; 
import axios from 'axios';
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Register', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {          
            Name: "",
            Uuid: "",
            Status:"",
           
      }
  },
  created() {

  },
  methods: {

    testMethod() {

    },
    async RegisterIdent() {
          console.log("RegisterIdent")
          if (this.Name === "") {
            this.$parent.showAlert("Kein Username angegeben!")
          } else {
            this.Uuid = uuid.v4()
            this.Status=""       
            //await   axios.post(apiURL + "/RegisterIdent", {  
              await   axios.post(apiURL + "/TestApi", {  
              MsgType:"RegisterRequest",         
              Name: this.Name,
              Uuid: this.Uuid
            }, {
              headers: {
                'Content-Type': 'application/json',
              }
            })
            .then (function(response) {
              console.log("Response: " + response.data.Error)
            })
            .catch((error)=> this.$parent.showAlert("Server returned an Error:\n" + error));  
    
          }
    },        
  }       
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
