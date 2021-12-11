<template>
  <div >    
    <p>
      Nacherfassen
    </p>   
        <button class="btn btn-info" v-on:click="TestBannerMessage()">Test-BannerMessage</button>
  </div>
</template>

<script>

export default {
  name: 'Nacherfassen',  
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
          this.$parent.showAlert("Nacherfassen")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
