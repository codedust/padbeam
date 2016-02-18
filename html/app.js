;(function() {
  function reloadPad() {
    $.get("/pad", function(data) {
      if (data == "RELOAD") {
        location.reload();
      } else {
        $("body").html(marked(data));
      }
    });
  }

  window.setInterval(reloadPad, 2*1000);
  reloadPad();
}());
