<template>
  <div
    ref="el"
    v-if="($route.path.startsWith('/new') || $route.path.startsWith('/edit')) && images.length > 0"
    :style="gridStyle"
    class="grid gap-2"
  >
    <div
      v-for="img in images"
      :key="img"
      class="relative"
    >
      <img
        :src="getImageUrl(img)"
        alt=""
        class="cursor-move rounded relative"
        :class="images.length === 1 ? 'full-cover-image-single' : 'full-cover-image-mult'"
      />
      <div
        class="absolute right-6 top-0 px-1 bg-white m-2 rounded hover:text-red-500 cursor-pointer"
        @click="removeImage(img)"
      >
        <UIcon name="i-carbon-trash-can" class="" />
      </div>
    </div>
  </div>

  <ViewImage v-else :style="gridStyle" v-if="images.length > 0">
    <img
      v-for="(img, z) in images"
      :key="z"
      class="cursor-zoom-in rounded"
      :class="images.length === 1 ? 'full-cover-image-single' : 'full-cover-image-mult'"
      :src="getImageUrl(img)"
      alt=""
    />
  </ViewImage>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useSortable } from '@vueuse/integrations/useSortable';
import type { SysConfigVO } from '~/types';
import ViewImage from './ViewImage.vue';

const sysConfig = useState<SysConfigVO>('sysConfig');
const route = useRoute();
const el = ref(null);
const props = defineProps<{ imgs: string }>();
const emit = defineEmits(['removeImage', 'dragImage']);

// Initialize images from the props
const images = ref<string[]>((!props.imgs || props.imgs === ',') ? [] : props.imgs.split(','));

// Watch for changes in props and update images accordingly
watch(props, () => {
  images.value = (!props.imgs || props.imgs === ',') ? [] : props.imgs.split(',');
});

// Function to get the image URL
const getImageUrl = (src: string) => {
  console.log(sysConfig.value.s3.thumbnailSuffix, src);
  if (src.startsWith('/')) {
    return src;
  }
  if (sysConfig.value.s3?.thumbnailSuffix) {
    const suffix = sysConfig.value.s3.thumbnailSuffix;
    if (src.indexOf(suffix) > 0) {
      return src;
    }
    return suffix.startsWith('?') ? `${src}${suffix}` : `${src}?${suffix}`;
  }
  return src;
};

// Watch images and emit 'dragImage' event
watch(images, () => {
  emit('dragImage', images.value);
});

// Function to remove an image
const removeImage = async (img: string) => {
  await useMyFetch('/memo/removeImage', { img });
  emit('removeImage', img);
};

// Initialize sortable functionality on mount
onMounted(() => {
  if (route.path.startsWith('/new') || route.path.startsWith('/edit')) {
    setTimeout(() => {
      useSortable(el, images);
    }, 500);
  }
});

// Compute grid style based on the number of images
const gridStyle = computed(() => {
  let style = 'max-width:100%;display:grid;gap: 0.5rem;align-items: start;';
  switch (images.value.length) {
    case 1:
      style += 'grid-template-columns: 1fr; max-width: 60%;';
      break;
    case 2:
      style += 'grid-template-columns: 1fr 1fr; aspect-ratio: 2 / 1;';
      break;
    case 3:
      style += 'grid-template-columns: 1fr 1fr 1fr; aspect-ratio: 3 / 1;';
      break;
    case 4:
      style += 'grid-template-columns: 1fr 1fr; aspect-ratio: 1;';
      break;
    default:
      style += 'grid-template-columns: 1fr 1fr 1fr;';
  }
  return style;
});
</script>

<style scoped>
.full-cover-image-mult {
  object-fit: cover;
  object-position: center;
  max-height: 300px;
  width: 100%;
  aspect-ratio: 1 / 1;
  border: transparent 1px solid;
}

.full-cover-image-single {
  object-fit: cover;
  object-position: center;
  max-height: 300px;
  height: auto;
  width: auto;
  border: transparent 1px solid;
}
</style>