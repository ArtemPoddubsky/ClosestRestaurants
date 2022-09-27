# Closest Restaurants 🍖
Service that provides API and HTML UI for interacting with data about Moscow restaurants.

<h1>Задачи</h1>

1. Загрузить <a href="https://github.com/ArtemPoddubsky/ClosestRestaurants/blob/main/Postgres/data/data.csv">data.csv</a> файл, спаршенный с <a href="https://data.mos.ru/">data.mos.ru</a> содержащий информацию о точках общественного питания в городе Москва, в базу данных.
    Каждая запись имеет поля: Id, Name, Address, Phone, Longitude, Latitude.
  
2. Создать HTML UI, используя шаблон, описанный в <a href="https://github.com/ArtemPoddubsky/ClosestRestaurants/blob/main/ClosestRestaurants/materials/page.html">page.html</a> . При запросе "localhost:5000/?page=2", ответ будет содержать сгенерированную html страницу, содержащую информацию об объектах базы данных, и простой интерфейс перелистывания страниц.

<p align="center">
    <image src="https://user-images.githubusercontent.com/108487635/188926812-9da4f892-ce9d-4be1-bd4c-560b3a908617.png">
</p>

  
  <h5 align=center>Пример ответа:</h5>
  
  ```
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Places</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>

<body>
<h5>Total: 13649</h5>
<ul>
    <li>
        <div>Sushi Wok</div>
        <div>gorod Moskva, prospekt Andropova, dom 30</div>
        <div>(499) 754-44-44</div>
    </li>
    <li>
        <div>Ryba i mjaso na ugljah</div>
        <div>gorod Moskva, prospekt Andropova, dom 35A</div>
        <div>(499) 612-82-69</div>
    </li>
    <li>
        <div>Hleb nasuschnyj</div>
        <div>gorod Moskva, ulitsa Arbat, dom 6/2</div>
        <div>(495) 984-91-82</div>
    </li>
    ...
</ul>
<a href="/?page=1">Previous</a>
<a href="/?page=3">Next</a>
<a href="/?page=1364">Last</a>
</body>
</html>
```

3. Написать API метод recommend, доступ к которому осушествляется по адресу "localhost:5000/api/recommend" с указанием координат (float) в теле запроса в формате json, по которым будет осуществлён поиск 3 ближайших локаций из базы данных.
    
<h5 align=center>Пример тела запроса:</h5>
    
```
{"lon":37.666, "lat":55.674}
```
    
<h5 align=center>Пример ответа:</h5>
    
```
{
  "name": "Recommendation",
  "places": [
    {
      "id": 30,
      "name": "Ryba i mjaso na ugljah",
      "address": "gorod Moskva, prospekt Andropova, dom 35A",
      "phone": "(499) 612-82-69",
      "location": {
        "lat": 55.67396575768212,
        "lon": 37.66626689310591
      }
    },
    {
      "id": 3348,
      "name": "Pizzamento",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.673075576456,
        "lon": 37.664533747576
      }
    },
    {
      "id": 3347,
      "name": "KOFEJNJa «KAPUChINOFF»",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.672865251005106,
        "lon": 37.6645689561318
      }
    }
  ]
}
```

<h2>Коды ответов</h2>

<h4>Global</h4>

    500 Internal Server Error: "Temporary Error (500)"

<h4>/?page=</h4>
    
    400 Bad Request: 
        page < 0: "This page can't possibly exist"
    
    404 Not Found:
        page > total pages: "This page doesn't exist"

<h4>/api/recommend</h4>

    400 Bad Request:
        No JSON, JSON syntax error, lat or lon omitted or incorrect type: "{"error": "Bad JSON"}"

<h2>Сборка</h2>

    make
    
    make lint – запуск линтера
