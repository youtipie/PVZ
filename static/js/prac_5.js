$(document).ready(function(){
    var data; // Відповідає за дані з /prac-5/data
    var element_count = 0; // Відповідає за те, скільки елементів ЕПС налічує сторінка
    // Здійснюємо ajax запит, щоб отримати перелік доступних елементів ЕПС
    $.ajax({
        url: '/prac-5/data',
        type: 'GET',
        success: function(response) {
            data = response;
        },
        error: function(xhr, status, error) {
            console.error('Error:', status, error);
        }
    });
    // Обробник події для видалення елемента ЕПС
    $('#dynamic-inputs').on('click', '.fa-delete-left', function(){
        // Видаляємо батьківський елемент, коли клікнуто на іконку видалення
        $(this).parent().remove();
        element_count -= 1;
        updateButtonState();
    });
    // Обробник події для додавання нового елемента ЕПС
    $('.fa-plus').on('click', function(){
        var newItem = $('<div class="input-group input-group-sm mt-3 mb-3">'+
                    '<label class="input-group-text fs-4 me-2">Кількість, елемент</label>'+
                    '<input type="number" name="quantity[]" value="1" min="1" step="1" class="form-control" required>'+
                    '<select name="element[]" class="form-select" required>'+
                        '<option value="">Оберіть елемент</option>'+
                    '</select>'+
                    '<i class="fa-solid fa-delete-left fa-2xl ms-4 mt-4" style="color: #d41616;"></i>'+
                '</div>');
        $('#dynamic-inputs').append(newItem); // Додаємо новий елемент до контейнера
        var selectElement = newItem.find('select');
        // Додаємо опції для вибору елемента з отриманих даних
        $.each(data, function(index, value) {
            selectElement.append($('<option>', {
                value: value,
                text: value
            }));
        });
        element_count += 1;
        updateButtonState();
    });
    // Функція для оновлення стану кнопки в залежності від кількості елементів
    // Зроблено для того, щоб користувач не міг відправити форму з пустими значеннями
    function updateButtonState() {
        if (element_count > 0){
            $('.btn-success').removeAttr('disabled');
        } else {
            $('.btn-success').attr('disabled','disabled');
        }
    }
});