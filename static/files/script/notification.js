document.addEventListener('DOMContentLoaded', () => {
  const notifelem = document.getElementById("notification");
  if (notifelem === null) {
    // no notification present
    return;
  }

  notifelem.addEventListener("click", () => {
    notifelem.remove();
  })
})
