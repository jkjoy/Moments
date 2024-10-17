<template>
  <div ref="container">
    <slot></slot>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

// 定义属性，以便自定义选择器和选项
const props = defineProps({
  selector: {
    type: String,
    default: '#view img, #view-2 img' 
  },
  options: {
    type: Object,
    default: () => ({}) // 默认选项为空对象
  }
});

// 对容器元素的引用
const container = ref(null);

// 生命周期钩子：在组件挂载时初始化 ViewImage
onMounted(() => {
  console.log('view-image initializing...');
  if (window.ViewImage) {
    window.ViewImage.init(props.selector, props.options); // 使用传入的选择器
    console.log('view-image initialized.');
  } else {
    console.error('ViewImage is not defined. Please ensure that view-image.min.js is loaded.');
  }
});

// 生命周期钩子：在组件卸载时销毁 ViewImage
onUnmounted(() => {
  if (window.ViewImage && window.ViewImage.destroy) {
    window.ViewImage.destroy();
    console.log('view-image destroyed.');
  }
});
</script>

<style scoped>
/* 在此添加任何需要的样式 */
</style>