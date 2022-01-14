<template>
  <div >    
    <p>
      Auswertung
    </p>   
  <table class="table table-dark">
  <thead>
    <tr>            
      <th scope="col">Von</th>
      <th scope="col">Bis</th>
    </tr>
  </thead>
  <tbody>
    <tr>   
      <td><input type="date" id="GetAccountingsFromDate" v-model="GetAccountingsFromDate" /> </td>
      <td><input type="date" id="GetAccountingsToDate" v-model="GetAccountingsToDate"/></td>      
    </tr>
  </tbody>
   <!-- {{this.GetAccountingsFromDate}} --> 
           <button class="btn btn-info" v-on:click="GetAccountings()">Anzeigen</button>
        <button class="btn btn-info" v-on:click="TestBannerMessage()">CSV-Download</button>
</table>

    <br>     
    <br>
               <table class="table table-dark">
              <thead>
                <tr>            
                  <th scope="col">Von</th>
                  <th scope="col">bis</th>
                  
                </tr>
              </thead>
              <tbody>
                <tr v-for="entry in this.Accountings" :key="entry.Id">   
                  <td>{{ entry.FromDate }}</td>
                  <td>{{ entry.ToDate }}</td>      
                  
                </tr>
              </tbody>            
            </table>

  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Auswertung',  
  props: {
    
  },
  data() {
      return {
          sendMessage: {
            MsgType: "",            
            Name: "",
            Typ:"",
            FromDate:""            
            },
          Accountings: [{}],
          GetAccountingsFromDate:"",
          GetAccountingsToDate:""
      }
  },
  created() {
    this.sendMessage.Name = this.$parent.username  

  },
  methods: {
      TestBannerMessage() {
          this.$parent.showAlert("CSV-Download")
      },
      async GetAccountings() {
       console.log("GetAccountings: " + this.$parent.APIURL + '/TestApi')
        this.sendMessage.MsgType = "GetAccountings"
        this.sendMessage.FromDate = this.GetAccountingsFromDate
        this.sendMessage.ToDate = this.GetAccountingsToDate

        axios.post (this.$parent.APIURL + '/TestApi', this.sendMessage)        
        .then ((res) => {
          console.log("Request: " + this.sendMessage)
          console.log("Result: " + JSON.stringify(res))
          this.Accountings = res.data
          if ( res.data === "") {
          this.$parent.showAlert("nothing found  Invalid result")  
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
