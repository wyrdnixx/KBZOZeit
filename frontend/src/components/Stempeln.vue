<template>
  <div >    
    <p>
      Stempeln
      Debu: {{this.sendMessage}}
    </p>   
        <button class="btn btn-info" v-on:click="TestButton()">Einstempeln</button>

        <br>
        <button class="btn btn-info" v-on:click="GetOpenTimer()">Test-offene Timer abfragen</button>
        <br>
        <div>
          <h1 v-if="this.openTimerStartTime == '' ">no open Timer</h1>
          <h1 v-else> OpenTimer:</h1>  
          
          <span>Start: {{ this.openTimerStartTime | moment("DD.MM.YYYY hh:mm:ss") }}</span>
          <br>
          <span>Now : {{ (new Date()) | moment("DD.MM.YYYY hh:mm:ss") }}</span>                    
          <br>
          <span>diff in hours:  {{this.TimeDiff}}</span>
          <br>

        </div>
        <button class="btn btn-info" v-on:click="getTimeDiff()">getTimeDiff</button>

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
            },
          openTimerStartTime:"",
          TimeDiff:""
          
      }
  },
  created() {
    this.sendMessage.Name = this.$parent.username  
    this.GetOpenTimer()

      //refresh - evtl noch verschieben - nur wenn offener timer gefunden wurde
      window.setInterval(() => {
        this.getTimeDiff()
      }, 10000) // every 10 Seconds
  },
  computed: {    
  },
  methods: {
    getTimeDiff(){
      //var diff =( (new Date(this.openTimerStartTime)) -(new Date())) / 1000;
      var now = new Date()
      var then = new Date(this.openTimerStartTime)
      console.log("now: " + now)
      console.log("then: " + then)
      var diff =( now - then)  / 1000;
      console.log("diff " + diff)
      diff /= (60 * 60);  // dif in stunden umgerechnet
      console.log("diff " + diff)
      console.log("diff " + diff.toFixed(2) )     

      this.TimeDiff =diff.toFixed(2) 
    },
      TestButton() {
          //this.$parent.showAlert("Stempeln")
          //this.TimeAccounting()
          this.TimeAccounting()
          
      },
      async GetOpenTimeaccounting(){
          console.log("GetOpenTimeaccounting")
          this.sendMessage.MsgType= "GetOpenTimeaccounting"
          this.sendMessage.Name= "testuser"
          axios.post(apiURL+'/TestApi',this.sendMessage)
          .then((res) => {
            console.log("Result: "+ res.data.Result)
          })
          .catch((error) => {                     
               console.log("Error:"+ error.response.data.Result)
               this.$parent.showAlert(error.response.data)
          })
           .finally(() => {
               //Perform action in always
           });
      },
      async GetOpenTimer() {
        console.log("TestButtonGetOpenTimer")
        this.sendMessage.Typ = "getOpenTimer"

        axios.post (apiURL + '/TestApi', this.sendMessage)
        .then ((res) => {
          console.log("Result: " + res.data.FromDate.String)
          this.openTimerStartTime = res.data.FromDate.String
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
                     this.GetOpenTimer()
                 });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
