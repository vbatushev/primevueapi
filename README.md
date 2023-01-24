
# PriveVue API Scraper

Как известно, у [PrimeVue](https://www.primefaces.org/primevue) есть [SASS API](https://www.primefaces.org/designer/api/primevue/3.9.0/).

Данная утилита забирает данные этой страницы и на их основе создает два файла:

* _variables.scss
* variables.json

### Параметры при запуске

*-ver* — версия API, по умолчанию, равна "3.9.0"

### Структура JSON

Это массив секций. Каждая секция имеет заголовок (title) и массив элементов секции (items).

Элемент секции имеет три параметра:

* property — название SASS-переменной
* value — значение SASS-переменной
* comment — комментарий к данной переменной

> Утилита создана исключительно для моего удобства