<template>
  <div>
    <div v-for="comment in commentArray" :key="comment.postId">
      <nested-comment :comment="comment" v-on:refresh-event="$emit('refresh-event')"/>
    </div>
    
  </div>
</template>

<script>
import NestedComment from "@/components/NestedComment.vue"
export default {
  data() {
    return {
    };

  },
  props: {
    commentArray: Array
  },
  components: {
    NestedComment
  },
};
</script>

<style>

</style>