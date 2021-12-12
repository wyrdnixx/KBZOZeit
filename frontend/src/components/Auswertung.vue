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
      <td><input type="date" id="DateFrom" /></td>
      <td><input type="date" id="DateTo" /></td>      
    </tr>
  </tbody>
</table>
        <button class="btn btn-info" v-on:click="TestBannerMessage()">Anzeigen</button>
        <button class="btn btn-info" v-on:click="TestBannerMessage()">CSV-Download</button>
    <br> <hr>
    <table class="table table-dark">
      <thead>
        <tr>            
          <th scope="col">Ergebnis</th>
          <th scope="col">Ergebnis</th>
        </tr>
      </thead>
      <tbody>
        <tr>   
          <td>Eintrag</td>      
          <td>Eintrag</td>
        </tr>
      </tbody>
    </table>

  </div>
</template>

<script>

export default {
  name: 'Auswertung',  
  props: {
    
  },
  data() {
      return {

      }
  },
  created() {

  },
  methods: {
      TestBannerMessage() {
          this.$parent.showAlert("CSV-Download")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
