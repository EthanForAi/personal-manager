const actions = {
  create: {
    endpoint: "/create",
    button: "新建资料",
    title: "已准备新建",
    copy: "提交后会生成一条新的个人资料。",
    needsFullPerson: true,
  },
  read: {
    endpoint: "/read",
    button: "查询资料",
    title: "已准备查询",
    copy: "根据 UserID 查询一条个人资料。",
    needsFullPerson: false,
  },
  update: {
    endpoint: "/update",
    button: "更新资料",
    title: "已准备更新",
    copy: "提交后会覆盖当前 UserID 对应的资料。",
    needsFullPerson: true,
  },
  delete: {
    endpoint: "/delete",
    button: "删除资料",
    title: "已准备删除",
    copy: "提交后会删除当前 UserID 对应的资料。",
    needsFullPerson: false,
  },
  check: {
    endpoint: "/check",
    button: "检查存在",
    title: "已准备检查",
    copy: "返回当前 UserID 是否已经存在。",
    needsFullPerson: false,
  },
};

let currentAction = "create";
let selectedUserID = "alice01";
let records = [
  {
    userid: "alice01",
    name: "Alice",
    email: "alice@example.com",
    phone: "13800138000",
  },
  {
    userid: "bob02",
    name: "Bob",
    email: "bob@example.com",
    phone: "13900139000",
  },
  {
    userid: "chen03",
    name: "Chen",
    email: "chen@example.com",
    phone: "18800188000",
  },
];

const form = document.querySelector("#person-form");
const actionButtons = Array.from(document.querySelectorAll("[data-action]"));
const fullFieldGroups = Array.from(document.querySelectorAll("[data-full-fields]"));
const activeEndpoint = document.querySelector("#active-endpoint");
const submitAction = document.querySelector("#submit-action");
const resetAction = document.querySelector("#reset-action");
const requestPreview = document.querySelector("#request-preview");
const responsePreview = document.querySelector("#response-preview");
const statusBadge = document.querySelector("#status-badge");
const statusCard = document.querySelector("#status-card");
const statusTitle = document.querySelector("#status-title");
const statusCopy = document.querySelector("#status-copy");
const recordList = document.querySelector("#record-list");
const recordCount = document.querySelector("#record-count");

function formValue(name) {
  return form.elements[name].value.trim();
}

function payloadForAction() {
  const payload = {
    userid: formValue("userid"),
  };

  if (actions[currentAction].needsFullPerson) {
    payload.name = formValue("name");
    payload.email = formValue("email");
    payload.phone = formValue("phone");
  }

  return payload;
}

function prettyJSON(value) {
  return JSON.stringify(value, null, 2);
}

function renderRequestPreview() {
  requestPreview.textContent = `POST ${actions[currentAction].endpoint}\n\n${prettyJSON(payloadForAction())}`;
}

function setStatus({ ok, badge, title, copy, response }) {
  statusBadge.textContent = badge;
  statusBadge.className = ok ? "badge success" : "badge danger";
  statusCard.classList.toggle("is-error", !ok);
  statusTitle.textContent = title;
  statusCopy.textContent = copy;
  responsePreview.textContent = prettyJSON(response);
}

function setAction(action) {
  currentAction = action;
  const meta = actions[action];

  actionButtons.forEach((button) => {
    button.classList.toggle("is-active", button.dataset.action === action);
  });

  fullFieldGroups.forEach((group) => {
    group.hidden = !meta.needsFullPerson;
  });

  activeEndpoint.textContent = meta.endpoint;
  submitAction.textContent = meta.button;
  setStatus({
    ok: true,
    badge: "Ready",
    title: meta.title,
    copy: meta.copy,
    response: {
      status: "prototype",
      endpoint: meta.endpoint,
    },
  });
  renderRequestPreview();
}

function validatePayload(payload) {
  if (!payload.userid) {
    return "userid is required";
  }

  if (payload.userid.length > 32) {
    return "userid must be no more than 32 characters";
  }

  if (!/^[A-Za-z][A-Za-z0-9]*$/.test(payload.userid)) {
    return "userid must start with a letter and contain letters and digits only";
  }

  if (!actions[currentAction].needsFullPerson) {
    return "";
  }

  if (!payload.name) {
    return "name is required";
  }

  if (/[0-9]/.test(payload.name)) {
    return "name must not contain digits";
  }

  if (!payload.email) {
    return "email is required";
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(payload.email)) {
    return "email must be a valid email address";
  }

  if (!payload.phone) {
    return "phone is required";
  }

  if (!/^1[3-9][0-9]{9}$/.test(payload.phone)) {
    return "phone must be a valid mainland China mobile number";
  }

  return "";
}

function findRecord(userid) {
  return records.find((record) => record.userid === userid);
}

function handlePrototypeSubmit(event) {
  event.preventDefault();

  const payload = payloadForAction();
  const validationError = validatePayload(payload);

  if (validationError) {
    setStatus({
      ok: false,
      badge: "400",
      title: "校验失败",
      copy: validationError,
      response: {
        error: validationError,
      },
    });
    return;
  }

  const existing = findRecord(payload.userid);

  if (currentAction === "create") {
    if (existing) {
      setStatus({
        ok: false,
        badge: "400",
        title: "UserID 已存在",
        copy: "这条记录已经在草稿列表中。",
        response: {
          error: "userid already exists",
        },
      });
      return;
    }

    records = [payload, ...records];
    selectedUserID = payload.userid;
    renderRecords();
    setStatus({
      ok: true,
      badge: "200",
      title: "新建成功",
      copy: "草稿列表已经加入这条资料。",
      response: payload,
    });
    return;
  }

  if (currentAction === "read") {
    if (!existing) {
      respondNotFound(payload.userid);
      return;
    }

    selectedUserID = existing.userid;
    fillForm(existing);
    renderRecords();
    setStatus({
      ok: true,
      badge: "200",
      title: "查询成功",
      copy: "已把查询结果填回表单。",
      response: existing,
    });
    return;
  }

  if (currentAction === "update") {
    if (!existing) {
      respondNotFound(payload.userid);
      return;
    }

    records = records.map((record) => (record.userid === payload.userid ? payload : record));
    selectedUserID = payload.userid;
    renderRecords();
    setStatus({
      ok: true,
      badge: "200",
      title: "更新成功",
      copy: "草稿列表已经同步当前资料。",
      response: payload,
    });
    return;
  }

  if (currentAction === "delete") {
    if (!existing) {
      respondNotFound(payload.userid);
      return;
    }

    records = records.filter((record) => record.userid !== payload.userid);
    selectedUserID = records[0]?.userid || "";
    renderRecords();
    clearForm();
    setStatus({
      ok: true,
      badge: "200",
      title: "删除成功",
      copy: "草稿列表已经移除这条资料。",
      response: {
        deleted: true,
      },
    });
    return;
  }

  if (currentAction === "check") {
    setStatus({
      ok: true,
      badge: "200",
      title: existing ? "记录存在" : "记录不存在",
      copy: existing ? "当前 UserID 已在草稿列表中。" : "当前 UserID 可以用于新建资料。",
      response: {
        exists: Boolean(existing),
      },
    });
  }
}

function respondNotFound(userid) {
  setStatus({
    ok: false,
    badge: "404",
    title: "没有找到记录",
    copy: `${userid} 不在草稿列表中。`,
    response: {
      error: "record not found",
    },
  });
}

function fillForm(record) {
  form.elements.userid.value = record.userid;
  form.elements.name.value = record.name;
  form.elements.email.value = record.email;
  form.elements.phone.value = record.phone;
  renderRequestPreview();
}

function clearForm() {
  form.reset();
  renderRequestPreview();
}

function renderRecords() {
  recordCount.textContent = `${records.length} records`;

  if (records.length === 0) {
    recordList.innerHTML = '<p class="pm-empty">暂无草稿记录</p>';
    return;
  }

  recordList.innerHTML = records
    .map((record) => {
      const selectedClass = record.userid === selectedUserID ? " is-selected" : "";
      return `
        <button type="button" class="pm-record-card${selectedClass}" data-userid="${record.userid}">
          <h3>
            <span>${record.name}</span>
            <code>${record.userid}</code>
          </h3>
          <dl>
            <dt>Email</dt>
            <dd>${record.email}</dd>
            <dt>Phone</dt>
            <dd>${record.phone}</dd>
          </dl>
        </button>
      `;
    })
    .join("");
}

actionButtons.forEach((button) => {
  button.addEventListener("click", () => setAction(button.dataset.action));
});

form.addEventListener("input", renderRequestPreview);
form.addEventListener("submit", handlePrototypeSubmit);
resetAction.addEventListener("click", clearForm);

recordList.addEventListener("click", (event) => {
  const card = event.target.closest("[data-userid]");
  if (!card) {
    return;
  }

  const record = findRecord(card.dataset.userid);
  if (!record) {
    return;
  }

  selectedUserID = record.userid;
  fillForm(record);
  renderRecords();
});

renderRecords();
fillForm(records[0]);
setAction("create");
