# Глоссарий перевода EN → RU

> Обязателен для всех, кто переводит Mathua: интерфейс, данные, уроки, генераторы задач.
> Если нужного термина нет, выберите вариант, добавьте его сюда в своём коммите и упомяните это в отчёте.

## 1. Правила стиля

- **Обращение — на «вы»**, со строчной буквы. Задания формулируются как в учебнике: «Найдите…», «Вычислите…», «Упростите…».
- **Регистр.** С заглавной пишется только первое слово кнопки, заголовка или пункта меню: «Пройти диагностику», а не «Пройти Диагностику».
- **Типографика.** Кавычки «ёлочки», внутри них „лапки“. Тире «—» с пробелами. Буква «ё» пишется.
- **Склонение с числами** только через plural-формы (one/few/many): «1 тема, 3 темы, 5 тем». Никаких «тем(ы)».
- **Длина.** Русский примерно на 20–30% длиннее английского. Для коротких кнопок и вкладок выбирайте самый короткий корректный вариант.

## 2. Что НЕЛЬЗЯ переводить или менять

- LaTeX целиком: `\( … \)`, `\[ … \]`, `$ … $`, `$$ … $$`, команды (`\frac`, `\tan`, `\sin`…). Внутри формул ничего не меняется: `\tan` не превращается в `\operatorname{tg}`, точка не превращается в запятую.
- Ответы (`Answer`), идентификаторы концептов (`frac.add.word`), ключи JSON, роуты (`/learn`), имена файлов и пути к картинкам, URL, код, плейсхолдеры (`%d`, `%s`, `{count}`, `{{name}}`).
- **Числа в условиях и ответах — с десятичной точкой** («2.5», а не «2,5»). Поле ввода и грейдер ждут точку. Запятую можно ставить только в чисто декоративном тексте интерфейса через `Intl.NumberFormat`, но не в задачах.
- Названия брендов и авторов: Mathua, Math Academy, Algebrica, Antonio Lupetti, KaTeX, MathLive.
- «XP» остаётся «XP» (при необходимости в пояснении: «очки опыта»).
- Единицы измерения и валюта в задачах не пересчитываются: «$5» → «5 долларов», метры остаются метрами.

## 3. Продукт и интерфейс

| EN | RU | Комментарий |
|---|---|---|
| Learn (раздел) | Учиться | пункт навигации |
| Review | Повторение | раздел; действие — «Повторить» |
| Reviews due | Пора повторить | |
| Graph (раздел, граф концептов) | Граф | НЕ «график» |
| Domains | Разделы | |
| Domain / subdomain | раздел / подраздел | |
| Leaderboard | Рейтинг | |
| League | лига | |
| Bronze / Silver / Gold / Platinum / Diamond | Бронза / Серебро / Золото / Платина / Алмаз | Emerald — Изумруд, Ruby — Рубин, Sapphire — Сапфир |
| Glossary | Глоссарий | |
| Contact | Связаться | |
| Profile | Профиль | |
| History | История | |
| Missed (фильтр) | Ошибки | |
| Settings | Настройки | |
| Admin | Администрирование | в узком меню — «Админ» |
| Log in / Login | Войти / Вход | |
| Sign up | Зарегистрироваться / Регистрация | |
| Log out | Выйти | |
| Guest | гость | |
| Account | аккаунт | |
| Share / share link | Поделиться / ссылка для просмотра | |
| Parent / teacher | родитель / учитель | |
| Concept (узел графа) | тема | в UI для ученика всегда «тема». В технической документации допустимо «концепт» |
| Prerequisite(s) | опорная тема / опорные темы | |
| Locked / unlocked | закрыта / открыта | о теме |
| Frontier | граница знаний | |
| Course | курс | |
| Lesson | урок | |
| Knowledge point (KP) | шаг урока | |
| Subgoal | подцель | |
| Worked example | разобранный пример | |
| Practice | практика / тренировка | |
| Problem / question | задача / вопрос | |
| Answer | ответ | |
| Hint | подсказка | |
| Explanation | объяснение / разбор | |
| Solution | решение | |
| Submit / Check | Проверить | кнопка отправки ответа |
| Skip | Пропустить | |
| Next / Continue | Дальше / Продолжить | |
| Retry / Try again | Попробовать ещё раз | |
| Correct! / Incorrect | Верно! / Неверно | |
| Diagnostic | диагностика | НЕ «тест» и НЕ «контрольная» |
| Take the diagnostic / Retake | Пройти диагностику / Пройти заново | |
| Recommended | Рекомендуем | |
| Quiz (50 XP gate) / mastery check | контрольная | «Take Test» → «Пройти контрольную» |
| Mastery / mastered | освоение / освоена | |
| Unseen / Learning / Practicing / Mastered / Decaying | Не начата / Изучается / Закрепляется / Освоена / Забывается | состояния темы, род женский («тема») |
| Streak (правильные подряд) | серия | «серия из 10 верных ответов» |
| Daily streak | дней подряд | |
| Daily goal | цель на день | |
| Weekly | за неделю / еженедельный | |
| Promote / demote (лига) | повышение / понижение | |
| Pause / paused | пауза / на паузе | |
| Accommodated timing | увеличенное время | |
| Spaced repetition | интервальное повторение | |
| Speed / accuracy | скорость / точность | |
| Time limit / threshold | лимит времени / порог | |
| Progress | прогресс | |
| Goals | Цели | |
| Plan | План | |
| Onboarding | знакомство / первые шаги | |
| Note | заметка | |
| Report (проблему) | Сообщить об ошибке | |
| Contributors / invitations | участники / приглашения | |
| Owner / admin / moderator | владелец / администратор / модератор | |
| Loading… | Загрузка… | |
| Something went wrong | Что-то пошло не так | |

## 4. Курсы

| EN | RU |
|---|---|
| 4th–8th Grade Math | Математика, 4 класс … Математика, 8 класс |
| Prealgebra | Введение в алгебру |
| Algebra I / Algebra II | Алгебра I / Алгебра II |
| Geometry | Геометрия |
| Precalculus | Алгебра и начала анализа |
| Calculus I / II / III | Математический анализ I / II / III |
| Linear Algebra | Линейная алгебра |
| Discrete Math / Discrete Math II | Дискретная математика / Дискретная математика II |
| Probability & Statistics | Теория вероятностей и статистика |
| Probability Theory | Теория вероятностей |
| Number Theory | Теория чисел |
| Differential Equations | Дифференциальные уравнения |
| Abstract Algebra / II | Общая алгебра / Общая алгебра II |
| Topology | Топология |
| Machine Learning | Машинное обучение |
| Complex Numbers and Analysis | Комплексные числа и комплексный анализ |
| Graph Theory | Теория графов |
| Numerical Analysis | Численные методы |
| Real Analysis | Вещественный анализ |

## 5. Разделы (domain id → RU)

| id | RU |
|---|---|
| arithmetic | Арифметика |
| fractions | Дроби |
| prealgebra | Введение в алгебру |
| algebra | Алгебра |
| geometry | Геометрия |
| trigonometry | Тригонометрия |
| precalculus | Начала анализа |
| calculus | Математический анализ |
| linear_algebra | Линейная алгебра |
| discrete_math | Дискретная математика |
| graph_theory | Теория графов |
| number_theory | Теория чисел |
| probability_theory | Теория вероятностей |
| statistics | Статистика |
| differential_equations | Дифференциальные уравнения |
| complex_numbers | Комплексные числа |
| abstract_algebra | Общая алгебра |
| real_analysis | Вещественный анализ |
| topology | Топология |
| numerical_analysis | Численные методы |
| machine_learning | Машинное обучение |

## 6. Математические термины

### Арифметика и дроби
| EN | RU |
|---|---|
| whole numbers | целые неотрицательные числа |
| natural numbers / integers | натуральные / целые числа |
| place value / digit | разряд / цифра |
| sum / difference / product / quotient | сумма / разность / произведение / частное |
| dividend / divisor / remainder | делимое / делитель / остаток |
| factor (множитель в произведении) | множитель |
| factor (делитель числа) | делитель |
| multiple | кратное |
| prime / composite | простое / составное число |
| GCD (GCF) / LCM | наибольший общий делитель (НОД) / наименьшее общее кратное (НОК) |
| fraction / numerator / denominator | дробь / числитель / знаменатель |
| proper / improper fraction | правильная / неправильная дробь |
| mixed number | смешанное число |
| lowest terms / simplest form | несократимая дробь; «в несократимом виде» |
| equivalent fractions | равные дроби |
| common denominator | общий знаменатель |
| reciprocal | обратное число |
| decimal | десятичная дробь |
| percent | процент |
| ratio / proportion / rate | отношение / пропорция / скорость (норма — по контексту) |
| round to the nearest … | округлите до … |
| estimate | оцените / прикидка |
| number line | числовая прямая |
| absolute value | модуль |
| order of operations | порядок действий |

### Алгебра
| EN | RU |
|---|---|
| expression / equation / inequality | выражение / уравнение / неравенство |
| variable / constant / coefficient | переменная / константа / коэффициент |
| term / like terms | слагаемое (член) / подобные слагаемые |
| monomial / binomial / polynomial | одночлен / двучлен / многочлен |
| degree (of polynomial) | степень |
| power / exponent / base | степень / показатель степени / основание |
| square root / radical | квадратный корень / корень (радикал) |
| simplify / expand / factor (verb) | упростите / раскройте скобки / разложите на множители |
| solve for x | решите относительно x / найдите x |
| linear / quadratic | линейное / квадратное (уравнение); квадратичная (функция) |
| slope / y-intercept | угловой коэффициент / точка пересечения с осью y |
| system of equations | система уравнений |
| function / domain / range | функция / область определения / множество значений |
| graph (of a function) | **график** (не «граф») |
| root / zero | корень / нуль функции |
| logarithm / exponential | логарифм / показательная (экспоненциальная) функция |
| sequence / series | последовательность / ряд |
| arithmetic / geometric sequence | арифметическая / геометрическая прогрессия |

### Геометрия и тригонометрия
| EN | RU |
|---|---|
| point / line / segment / ray | точка / прямая / отрезок / луч |
| angle: acute / right / obtuse / straight | угол: острый / прямой / тупой / развёрнутый |
| parallel / perpendicular | параллельные / перпендикулярные |
| triangle: equilateral / isosceles / scalene | треугольник: равносторонний / равнобедренный / разносторонний |
| hypotenuse / leg | гипотенуза / катет |
| congruent / similar | равные / подобные |
| perimeter / area / volume / surface area | периметр / площадь / объём / площадь поверхности |
| circle (линия) / disk (область) | окружность / круг |
| radius / diameter / circumference | радиус / диаметр / длина окружности |
| polygon / quadrilateral / rectangle / square | многоугольник / четырёхугольник / прямоугольник / квадрат |
| sine / cosine / tangent | синус / косинус / тангенс (в тексте); в LaTeX остаётся `\tan` |
| unit circle / radian / degree | единичная окружность / радиан / градус |

### Высшая математика
| EN | RU |
|---|---|
| limit / derivative / integral | предел / производная / интеграл |
| antiderivative | первообразная |
| converge / diverge | сходится / расходится |
| matrix / vector / determinant | матрица / вектор / определитель |
| eigenvalue / eigenvector | собственное значение / собственный вектор |
| vector space / basis / span | векторное пространство / базис / линейная оболочка |
| set / subset / element | множество / подмножество / элемент |
| group / ring / field | группа / кольцо / поле |
| probability / event / sample space | вероятность / событие / пространство элементарных исходов |
| expected value / variance / standard deviation | математическое ожидание / дисперсия / стандартное отклонение |
| mean / median / mode | среднее / медиана / мода |
| sample / population | выборка / генеральная совокупность |
| graph (теория графов) / vertex / edge | граф / вершина / ребро |
| proof / theorem / lemma / corollary | доказательство / теорема / лемма / следствие |
| real numbers | действительные числа |
