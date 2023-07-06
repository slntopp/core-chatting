<template>
    <div :style="{flexDirection: isRight?'row-reverse':'row'}" class="mock-message">
      <n-skeleton  class="avatar" circle/>
      <div class="text">
        <div  class="text-header" :style="{flexDirection:isRight?'row-reverse':'row',justifyContent:isRight?'end':'start'}">
          <n-skeleton round :width="widthOfTitle+'px'" height="40px" text/>
          <n-skeleton round class="title" width="60px" height="25px" text/>
        </div>
        <n-skeleton width="100%" :height="heightOfText+'px'" text></n-skeleton>
      </div>
    </div>
</template>
<script setup lang="ts">
import {NSkeleton} from 'naive-ui'
import {onMounted, ref} from "vue";

const isRight=ref(false);
const heightOfText=ref(30);
const widthOfTitle=ref(80)


onMounted(()=>{

  heightOfText.value=randomIntFromInterval(1,5)*25
  widthOfTitle.value=randomIntFromInterval(70,200)
  isRight.value=randomIntFromInterval(1,2)===1
})

function randomIntFromInterval(min, max) { // min and max included
  return Math.floor(Math.random() * (max - min + 1) + min)
}
</script>

<style scoped lang="scss">
  .mock-message{
    width: 100%;
    display: flex;
    padding: 12px 0 12px 12px;
    .avatar{
      width: 72px;
      height: 64px;
      margin: 0px 15px;
    }
    .text{
      width: 80%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      .text-header{
        margin-bottom: 15px;
        display: flex;
        justify-content: start;
        align-items: start;
        .title{
          margin: 0px 5px;
        }
      }
    }
  }
</style>