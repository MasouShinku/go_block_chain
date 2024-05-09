<template>
  <div class="card">
    <div class="card-header pb-0 px-3">
      <h6 class="mb-0">{{ title }}</h6>
    </div>
    <div class="card-body pt-4 p-3">
      <ul class="list-group">
        <li
          v-for="({ Hash, PreviousHash, Timestamp, Pow,Trades }, index) of bills"
          :key="index"
          class="list-group-item border-0 d-flex p-4 mb-2 bg-gray-100 border-radius-lg"
        >
          <div class="d-flex flex-column">
            <h6 class="mb-3 text-sm">{{ Hash }}</h6>
            <span class="mb-2 text-xs">
              前向哈希:
              <span class="text-dark font-weight-bold ms-sm-2">
                {{ PreviousHash }}</span
              >
            </span>
            <span class="mb-2 text-xs">
              时间戳:
              <span class="text-dark ms-sm-2 font-weight-bold">
                {{ Timestamp }}</span
              >
            </span>
            <span class="mb-2 text-xs">
              交易数:
              <span class="text-dark ms-sm-2 font-weight-bold">
                <pre>{{ Trades }}</pre></span
              >
            </span>
            <span class="text-xs">
              Pow验证:
              <span class="text-dark ms-sm-2 font-weight-bold">{{ Pow }}</span>
            </span>
          </div>
          <div class="ms-auto text-end">
            <!-- <a
              class="btn btn-link text-danger text-gradient px-3 mb-0"
              href="javascript:;"
            >
              <i class="far fa-trash-alt me-2" aria-hidden="true"></i>Delete
            </a>
            <a class="btn btn-link text-dark px-3 mb-0" href="javascript:;">
              <i class="fas fa-pencil-alt text-dark me-2" aria-hidden="true"></i
              >Edit
            </a> -->
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: "BillingCard",
  props: {
    title: {
      type: String,
      default: "",
    },
    bills: {
      type: Array,
      name: String,
      company: String,
      email: String,
      id: String,
      Hash:String,
      PreviousHash:String,
      Timestamp:String,
      Pow:Boolean,
      Trades:Array,

      default: () => [],
    },
  },
};
</script>
