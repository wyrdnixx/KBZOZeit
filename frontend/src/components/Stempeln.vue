<template>
  <div >    
    <p>
      Stempeln
      Debu: {{this.sendMessage}}
    </p>   
        <button class="btn btn-info" v-on:click="TestButton()">Einstempeln</button>

        <br>
        <button class="btn btn-info" v-on:click="TestButtonGetOpenTimer()">Test-offene Timer abfragen</button>


  </div>
</template>

<script>
import axios from 'axios';
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Stempeln',  
  props: {
    
  },
  data() {
      return {
          sendMessage: {
            MsgType: "TimeAccounting",            
            Name: "",
            Typ:""
            }
      }
  },
  created() {
    this.sendMessage.Name = this.$parent.username          
  },
  methods: {
      TestButton() {
          //this.$parent.showAlert("Stempeln")
          this.TimeAccounting()
      },
      async TestButtonGetOpenTimer() {
        console.log("TestButtonGetOpenTimer")
        this.sendMessage.Typ = "getOpenTimer"

        axios.post (apiURL + '/TestApi', this.sendMessage)
        .then ((res) => {
          console.log("Result: " + res.data.FromDate.String)
        })
        .catch((error) => {                     
          console.log("Error:"+ error.response.data.Result)
          this.$parent.showAlert(error.response.data)
        })
        .finally(() => {
          //Perform action in always
        });
      },

      async TimeAccounting() {
          console.log("TimeAccounting")
            
          // start or stop 
          this.sendMessage.Typ = "startAccounting"

             axios.post(apiURL+ '/TestApi', this.sendMessage)
                 .then((res) => {
                     //Perform Success Action
                     console.log("Resut: "+ res.data.Result)
                     if (res.data.Result == "login successfully") {
                       this.$cookies.set("username",this.sendMessage.Name);
                       this.$parent.checkCookie()
                     }
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
