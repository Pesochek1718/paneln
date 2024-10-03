<div id="createServerModal" class="modal"><div class="modal-dialog">
<div class="modal-content">
<div class="modal-header">
<h3 class="modal-title">Создания сервера</h3> <a href="#close" title="{{.LANG.Modal.control.close}}" class="close" onclick="closeModal('createServer')">X</a></div><div class="modal-body">
<div id="createServerErr"></div>
<div id="createServerSingle">
<div style="margin: 20px 0 0 0;"><p>IP адрес</p><input id="createServerIP" type="text" class="inputtext" value="" placeholder="Только IPv4 адрес"></div>
<div style="margin: 20px 0 0 0;"><p>Время ожидания</p><input id="createServerWaitTime" type="text" class="inputtext" value="1800" placeholder="Только цифрами в часах"></div>
<div style="margin: 20px 0 0 0;"><p>Имя пользователя</p><input id="createServerUsername" type="text" class="inputtext" value="root"></div>
<div style="margin: 20px 0 0 0;"><p>Пароль пользователя</p><input id="createServerPassword" type="text" class="inputtext" value="1Htaht;bhfnjh"></div>
<div style="margin: 20px 0 0 0;"><p>Ключ SSH</p><textarea id="createServerKeySsh" class="inputtext" rows="10" cols="50"></textarea></div>
</div>

<div style="text-align: center">
<input type="button" value="Добавить" onclick="createServer()" id="createServerButton" class="sendbutton">
</div>
</div>
</div>
</div>

</div>

