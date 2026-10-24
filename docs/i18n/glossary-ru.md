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

- LaTeX целиком: `\( … \)`, `\[ … \]`, `$ … $`, `$$ … $$`, команды (`\frac`, `\tan`, `\sin`…). Внутри формул ничего не меняется: `\tan` не превращается в `\operatorname{tg}`, точка не превращается в запятую. **Исключение:** слова внутри `\text{…}` (и `\textbf`, `\textit`, `\textrm`, `\mbox`) — это проза и единицы, их переводим: `\text{cm}` → `\text{см}`; всё остальное в формуле остаётся как в оригинале, валидатор уроков это проверяет.
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
| Prerequisite(s) | необходимая тема / необходимые темы | в заголовках панелей — «Что изучить сначала»; «Обязательная» без существительного. Не «опорная»: «опорный» занят математикой (опорные дроби, опорный угол). Решено 2026-10-23 |
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
| whole numbers | целые неотрицательные числа; **в разделе дробей — натуральные числа** («Деление натурального числа на дробь»): там речь не о нуле, а строгий оборот делает текст тяжёлым |
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

## 7. Дополнения: названия тем и шаги уроков (данные)

Термины, выбранные при переводе `data/i18n/ru/concepts.json` и `data/lessons/kp/ru/`. Если на экране один и тот же термин, пишите его так.

### Правила для данных

- **Десятичная запятая и точка.** В числах из условий, ответов и примеров остаётся точка (`2.5`), как и в разделе 2. В названиях тем и шагах уроков (`kp/ru`) говорим «десятичная запятая». В текстах уроков о десятичных дробях знак называем «десятичной точкой», потому что на экране рядом стоит `3.7` (решено 2026-10-24, см. §13 и §15).
- **Названия тем** — существительные («Сложение дробей», «Деление столбиком»), а не команды. **Подцели шагов** — команды на «вы»: «Найдите слагаемые и сумму».
- **Формулы, индексы и идентификаторы** (`x^a`, `a_n`, `\\(n\\)`, `log_b(1)`, числа) в названиях и подцелях не меняются.

### Арифметика, дроби, введение в алгебру

| EN | RU |
|---|---|
| carry / borrow (в столбик) | перенос / занять десяток; «переход через разряд» |
| addend / minuend / subtrahend | слагаемое / уменьшаемое / вычитаемое |
| base-10 blocks | блоки по основанию 10 |
| expanded form | развёрнутая запись |
| times table | таблица умножения |
| long division | деление столбиком (чисел); деление уголком (многочленов) |
| factor tree | дерево множителей |
| word problem | текстовая задача |
| sign rules | правила знаков |
| benchmark fractions | опорные дроби |
| complex fraction | многоэтажная дробь |
| unit fraction | единичная дробь |
| scientific notation | стандартный вид числа |
| percent error / percent change | погрешность в процентах / изменение в процентах |
| markup / commission / sales tax | наценка / комиссионные / налог с продаж |
| unit rate / unit price | величина на единицу / цена за единицу |
| cross products | произведения крест-накрест |

### Алгебра, геометрия, тригонометрия, анализ

| EN | RU |
|---|---|
| slope-intercept form | уравнение прямой с угловым коэффициентом |
| standard form of a line | общее уравнение прямой |
| sign analysis (неравенства) | метод интервалов |
| synthetic division | схема Горнера |
| notable products | формулы сокращённого умножения |
| complete the square | выделение полного квадрата |
| extraneous root | посторонний корень |
| piecewise function | кусочно заданная функция |
| composite function | композиция функций (в шагах — «сложная функция» допустима) |
| dilation | гомотетия |
| midsegment | средняя линия |
| inscribed / central angle | вписанный / центральный угол |
| cyclic quadrilateral | вписанный четырёхугольник |
| law of sines / cosines | теорема синусов / косинусов |
| reciprocal trig functions (sec, csc, cot) | секанс, косеканс, котангенс (не «обратные функции»: обратные — это арк-функции) |
| dot / cross product | скалярное / векторное произведение |
| mean value theorem | теорема Лагранжа о среднем значении |
| squeeze theorem | теорема о промежуточной функции («о двух милиционерах») |
| improper integral | несобственный интеграл |
| line integral | криволинейный интеграл |
| sigma notation | знак суммы Σ |

### Дискретная математика, вероятность, анализ данных, ML

| EN | RU |
|---|---|
| pigeonhole principle | принцип Дирихле |
| proof by contradiction / contrapositive | доказательство от противного / через контрапозицию |
| inclusion-exclusion | формула включений и исключений |
| stars and bars | метод шаров и перегородок |
| derangement | беспорядок |
| random variable | случайная величина |
| expected value | математическое ожидание |
| law of large numbers / CLT | закон больших чисел / центральная предельная теорема |
| confidence interval / margin of error | доверительный интервал / предельная ошибка выборки |
| p-value / significance | p-значение / уровень значимости |
| accuracy / precision / recall | доля правильных ответов (accuracy) / точность (precision) / полнота (recall) |
| overfitting | переобучение |
| backpropagation | обратное распространение ошибки |
| one-hot encoding | унитарное кодирование (one-hot) |
| recognize (подцель) | распознайте (не «узнайте» — читается как «выясните») |

### Теория графов, численные методы, топология

| EN | RU |
|---|---|
| cut vertex / bridge | точка сочленения / мост |
| matching / perfect matching | паросочетание / совершенное паросочетание |
| vertex cover / independent set | вершинное покрытие / независимое множество |
| list coloring / choosability | списковая раскраска / выбираемость |
| treewidth / genus | древесная ширина / род графа |
| minor / subdivision | минор / подразбиение |
| outerplanar graph | внешнепланарный граф |
| tournament | турнир |
| max-flow min-cut / residual graph | максимальный поток и минимальный разрез / остаточный граф |
| topological sort (Kahn's algorithm) | топологическая сортировка (алгоритм Кана) |
| roundoff error / machine epsilon | ошибка округления / машинный эпсилон |
| catastrophic cancellation | катастрофическая потеря точности при вычитании |
| condition number / stability | число обусловленности / устойчивость |
| bracketing method / open method | метод с отделением корня / открытый метод |
| fixed-point iteration | метод простой итерации |
| false position (regula falsi) | метод хорд (regula falsi); «метод ложного положения» — только как пояснение |
| Horner's method | схема Горнера |
| divided differences | разделённые разности |
| Runge's phenomenon | феномен Рунге |
| composite rule / adaptive quadrature | составная квадратурная формула / адаптивная квадратура |
| diagonal dominance / sweep | диагональное преобладание / проход (по неизвестным) |
| partial pivoting / pivot | частичный выбор ведущего элемента / ведущий элемент |
| open ball / clopen | открытый шар / открыто-замкнутое множество |
| quotient topology | фактортопология |
| path-connected | линейно связный |
| wedge of circles | букет окружностей |
| deck transformation | автоморфизм накрытия (скольжение) |
| crosscap | кросс-кэп |
| Betti numbers / Euler characteristic | числа Бетти / эйлерова характеристика |

### Общая алгебра, комплексные числа и комплексный анализ

| EN | RU |
|---|---|
| coset / normal subgroup | смежный класс / нормальная подгруппа |
| quotient group / quotient ring | факторгруппа / факторкольцо |
| kernel / image | ядро / образ |
| abelianization / commutator subgroup | абелианизация / коммутант |
| free group / presentation | свободная группа / задание образующими и соотношениями |
| orbit-stabilizer theorem | теорема об орбите и стабилизаторе |
| Sylow subgroup | силовская подгруппа |
| nilpotent / solvable group | нильпотентная / разрешимая группа |
| composition series | композиционный ряд |
| class equation | уравнение классов |
| UFD / PID / Euclidean domain | факториальное кольцо / кольцо главных идеалов / евклидово кольцо |
| integral domain / zero divisor | область целостности / делитель нуля |
| prime / maximal ideal | простой / максимальный идеал |
| localization | локализация |
| splitting field | поле разложения |
| Noetherian ring | нётерово кольцо |
| tensor product / base change | тензорное произведение / замена базы |
| exact sequence | точная последовательность |
| principal value / branch / branch cut | главное значение / ветвь / разрез |
| pole / essential singularity | полюс / существенно особая точка |
| residue | вычет |
| Laurent series (principal part / regular part) | ряд Лорана (главная часть / правильная часть) |
| entire function | целая функция |
| Cauchy-Riemann equations | условия Коши — Римана |
| conformal map / Möbius transformation | конформное отображение / дробно-линейное преобразование (преобразование Мёбиуса) |
| cross-ratio | двойное отношение |
| locus | геометрическое место точек |
| Riemann sphere | сфера Римана |

### Дифференциальные уравнения, линейная алгебра

| EN | RU |
|---|---|
| ODE / PDE | ОДУ / УЧП (уравнение в частных производных) |
| separable equation | уравнение с разделяющимися переменными |
| integrating factor | интегрирующий множитель |
| exact equation | уравнение в полных дифференциалах |
| undetermined coefficients / variation of parameters | метод неопределённых коэффициентов / метод вариации постоянных |
| Wronskian | вронскиан |
| Cauchy problem / initial value problem | задача Коши (начальные условия); IVP — тоже «задача Коши» |
| boundary value problem | краевая задача |
| equilibrium / phase portrait | положение равновесия / фазовый портрет |
| saddle / node / spiral | седло / узел / фокус (спираль) |
| pitchfork / saddle-node / transcritical bifurcation | бифуркация вилка / седло-узел / транскритическая |
| Lyapunov function | функция Ляпунова |
| characteristics (PDE) | характеристики |
| Green's function | функция Грина |
| Sturm-Liouville problem | задача Штурма — Лиувилля |
| matrix exponential | матричная экспонента |
| Heaviside step / delta | функция Хевисайда / дельта-функция |
| impulse response / transfer function | импульсная характеристика / передаточная функция |
| dot product / cross product (в шагах) | скалярное произведение / векторное произведение |
| tip to tail | «конец к началу» |
| cofactor / minor / adjugate | алгебраическое дополнение / минор / присоединённая матрица |
| rank / nullity | ранг / дефект |
| column space / null space | пространство столбцов / нуль-пространство (ядро) |
| algebraic / geometric multiplicity | алгебраическая / геометрическая кратность |
| span / basis / change of basis | линейная оболочка / базис / замена базиса |
| orthogonal complement / projection | ортогональное дополнение / проекция |
| positive definite / indefinite | положительно определённая / знаконеопределённая |
| Sylvester's criterion | критерий Сильвестра |
| Cayley-Hamilton theorem | теорема Гамильтона — Кэли |
| Jordan normal form / minimal polynomial | жорданова нормальная форма / минимальный многочлен |

### Машинное обучение, теория чисел

| EN | RU |
|---|---|
| epoch / mini-batch | эпоха / мини-батч | решено 2026-10-23, так говорят в русских курсах ML |
| learning rate | скорость обучения |
| loss function | функция потерь |
| forward pass / backward pass | прямой проход / обратный проход |
| bias (смещение нейрона) | смещение |
| activation function | функция активации |
| overfit / underfit | переобучение / недообучение |
| dropout / batch normalization | дропаут / пакетная нормализация |
| momentum | импульс |
| actor / critic / advantage / baseline | актор / критик / преимущество / базовая линия |
| return (RL) | суммарная награда (не «возврат») |
| inference (этап применения модели) | этап предсказания (не «вывод») |
| TPR / FPR / AUC | TPR / FPR / AUC (площадь под ROC-кривой) |
| confusion matrix | матрица ошибок |
| F1 score | F1-мера |
| calibration (ECE) | калибровка (ожидаемая ошибка калибровки, ECE) |
| elbow method | метод локтя |
| dummy variable trap | ловушка фиктивных переменных |
| Gini impurity | неопределённость Джини |
| bagging / boosting | бэггинг / бустинг |
| residue / quadratic residue | вычет / квадратичный вычет |
| witness (Miller-Rabin) | свидетель |
| primitive root | первообразный корень |
| Euler's totient function | функция Эйлера |
| twin primes | простые числа-близнецы |
| convergent (цепной дроби) | подходящая дробь |
| class number | число классов |
| Gaussian integers | гауссовы целые числа |
| Pisano period | период Пизано |
| squarefree | свободное от квадратов |
| Euler's criterion | критерий Эйлера |

### Дискретная математика, теория вероятностей

| EN | RU |
|---|---|
| sigma-algebra / measurable | σ-алгебра / измеримый |
| power set | булеан |
| pmf / pdf / CDF | функция вероятности (pmf) / плотность распределения (pdf) / функция распределения (CDF) |
| prior / likelihood / posterior | априорная вероятность / правдоподобие / апостериорная вероятность |
| law of total probability / expectation / variance | формула полной вероятности / математического ожидания / дисперсии |
| tower property | свойство повторного ожидания |
| moment generating function / characteristic function | производящая функция моментов / характеристическая функция |
| order statistic | порядковая статистика |
| memoryless property | свойство отсутствия памяти |
| Poisson process | пуассоновский процесс |
| exponential family / sufficient statistic | экспоненциальное семейство / достаточная статистика |
| weak / strong law of large numbers | слабый / усиленный закон больших чисел |
| convergence in probability / in distribution / almost surely | сходимость по вероятности / по распределению / почти наверное |
| Borel-Cantelli lemmas | леммы Бореля — Кантелли |
| Slutsky's theorem / delta method | теорема Слуцкого / дельта-метод |
| Markov chain / transition matrix / stationary distribution | цепь Маркова / матрица переходных вероятностей / стационарное распределение |
| random walk / martingale / optional stopping | случайное блуждание / мартингал / теорема об остановке |
| recurrent / transient | возвратный / невозвратный |
| union bound | неравенство Буля |
| LOTUS (law of the unconscious statistician) | правило «бессознательного статистика» |

### Вещественный анализ, статистика

| EN | RU |
|---|---|
| upper bound / supremum / infimum | верхняя граница / супремум / инфимум |
| completeness axiom / Archimedean property | аксиома полноты / свойство Архимеда |
| nested intervals | вложенные отрезки |
| dense set / countable set | плотное множество / счётное множество |
| extended real line | расширенная числовая прямая |
| Cauchy sequence / Cauchy criterion | последовательность Коши (фундаментальная) / критерий Коши |
| subsequence / cluster point | подпоследовательность / частичный предел |
| divergence test | необходимый признак сходимости (расходимости) |
| comparison / ratio / root test | признак сравнения / Даламбера / Коши |
| absolute / conditional convergence | абсолютная / условная сходимость |
| rearrangement of a series | перестановка членов ряда |
| sequential criterion (Heine) | критерий предела по Гейне (через последовательности) |
| intermediate value theorem / extreme value theorem | теорема о промежуточном значении / теорема Вейерштрасса об экстремальных значениях |
| uniform continuity | равномерная непрерывность |
| compact / Heine-Borel | компактный / теорема Гейне — Бореля |
| removable discontinuity / jump | устранимый разрыв / скачок |
| Riemann sum, mesh, upper/lower sums | сумма Римана, диаметр (мелкость) разбиения, верхние/нижние суммы |
| uniform convergence | равномерная сходимость |
| contraction mapping | сжимающее отображение |
| nowhere dense / meagre | нигде не плотное / первой категории |
| Lipschitz | липшицев |
| mean, geometric, harmonic, quadratic mean | среднее арифметическое, геометрическое, гармоническое, квадратичное |
| Hölder mean | среднее Гёльдера |
| range / interquartile range | размах / межквартильный размах |
| z-score | z-оценка |
| standard error / margin of error | стандартная ошибка / предельная ошибка выборки |
| null / alternative hypothesis | нулевая / альтернативная гипотеза |
| Type I / Type II error / power | ошибка I рода / II рода / мощность |
| ANOVA | дисперсионный анализ |
| coefficient of determination (R-squared) | коэффициент детерминации (R-квадрат) |

### Начала анализа, тригонометрия

| EN | RU |
|---|---|
| terminal side / initial ray | конечная сторона / начальный луч |
| coterminal angles | углы с общей конечной стороной |
| quadrantal angle | осевой угол (лежащий на оси) |
| bearing | азимут |
| reduction formulas | формулы приведения |
| cofunction identities | тождества для дополнительных углов (кофункции) |
| double / half / triple angle formulas | формулы двойного / половинного / тройного угла |
| product-to-sum / sum-to-product | преобразование произведения в сумму / суммы в произведение |
| power-reduction formulas | формулы понижения степени |
| ambiguous case (SSA) | неоднозначный случай («сторона — сторона — угол») |
| Heron's formula / semiperimeter | формула Герона / полупериметр |
| amplitude / midline / phase shift | амплитуда / средняя линия / сдвиг по фазе |
| hyperbolic sine / cosine (sinh, cosh) | гиперболические синус / косинус |
| principal branch | главная ветвь |
| polar coordinates (pole, polar axis) | полярные координаты (полюс, полярная ось) |
| parametric equations | параметрические уравнения |
| directrix | директриса |
| partial fractions | простейшие дроби |
| piecewise function / sgn | кусочно заданная функция / функция знака sgn |
| floor / ceiling | «пол» (целая часть) / «потолок» |
| Heaviside step function | ступенчатая функция Хевисайда |
| arithmetic / geometric series (precalc) | арифметическая / геометрическая прогрессия, ряд |

### Математический анализ

| EN | RU |
|---|---|
| limit laws / squeeze theorem | свойства пределов / теорема о промежуточной функции (о зажатой функции) |
| epsilon-delta definition | определение на языке ε–δ |
| one-sided limit / infinite limit | односторонний предел / бесконечный предел |
| little-o / big-O | o-малое / O-большое |
| convex (f'' > 0) / concave (f'' < 0) | выпуклая вниз / выпуклая вверх (в шагах: «чаша вверх» / «чаша вниз») |
| inflection point | точка перегиба |
| cusp / corner / vertical tangent | точка возврата / угловая точка / вертикальная касательная |
| stationary point / critical point | стационарная / критическая точка |
| differential | дифференциал |
| related rates | связанные скорости изменения |
| linear approximation / linearization | линейное приближение / линеаризация |
| solid of revolution: disc / washer / shell method | тело вращения: метод дисков / шайб / цилиндрических оболочек |
| net change theorem | теорема о суммарном изменении |
| improper integral (comparison, p-test) | несобственный интеграл (признак сравнения, p-признак) |
| integration by parts / partial fractions / substitution | интегрирование по частям / разложением на простейшие дроби / заменой переменной |
| Weierstrass substitution | подстановка Вейерштрасса (тангенс половинного угла) |
| series: ratio / root / integral / Leibniz / Dirichlet test | ряд: признаки Даламбера / Коши (радикальный) / интегральный / Лейбница / Дирихле |
| telescoping series | телескопический ряд |
| Cauchy product | произведение рядов по Коши |
| power series, radius of convergence | степенной ряд, радиус сходимости |
| Taylor / Maclaurin / Fourier series | ряд Тейлора / Маклорена / Фурье |
| partial derivative / gradient / directional derivative | частная производная / градиент / производная по направлению |
| Jacobian | якобиан (матрица Якоби) |
| Hessian | гессиан |
| Lagrange multipliers | множители Лагранжа |
| double / triple integral | двойной / тройной интеграл |
| line integral (scalar) / surface integral | криволинейный / поверхностный интеграл (первого рода) |
| curl / divergence / flux | ротор / дивергенция / поток |
| Green / Stokes / divergence theorem | теорема Грина / Стокса / Остроградского — Гаусса |
| conservative field / potential | потенциальное поле / потенциал |
| curvature / torsion / Frenet frame | кривизна / кручение / репер Френе |
| gamma / beta function | гамма- / бета-функция |

## 8. Дополнения от перевода веб-интерфейса

> Добавлено при переводе `web/next-app` (2026-10-22). Термины выбраны там, где §3 не давал готового варианта.

### Интерфейс

| EN | RU | Комментарий |
|---|---|---|
| Next up | Что дальше | заголовок карточки следующего задания |
| Remediation (карточка рекомендации) | Восполнить пробел | повторение необходимой темы, которая блокирует следующие темы |
| Reference (панель внутри урока) | Справка | не «справочник»: это панель при теме |
| Reference (страница глоссария) | Справочник | метка раздела |
| Step N of M (шаг урока) | Шаг {n} из {total} | |
| Report a problem | Сообщить об ошибке | кнопка жалобы на вопрос |
| Closed book | Без подсказок | о контрольной |
| Practice again | Тренироваться снова | |
| Mastery check | Проверка освоения | заголовок экрана контрольной |
| Destination (в планировщике) | направление | |
| Rest days | выходные | |
| Workload | объём / нагрузка | по контексту |
| Provisional | предварительно | |
| Guest run / demo | гостевое прохождение / демо | |
| Strong (результат диагностики) | Уверенно | |
| Needs review | Нужно повторить | |
| Display name | отображаемое имя | |
| Username | имя пользователя | |
| Recovery email | почта для восстановления | |
| Connected accounts | подключённые аккаунты | |
| Pause learning | пауза в обучении | |
| Resets in (рейтинг) | До обнуления | |
| Learner (колонка рейтинга) | Ученик | |
| Attempt source: practice / quiz / diagnostic / review | практика / контрольная / диагностика / повторение | фильтр «Источник» |
| Creator's note | Заметка автора | пункт меню и заголовок страницы |
| Retention | Удержание | |
| Curriculum (метка) | Программа | |
| Title block (метка на лендинге) | Основная надпись | |
| Graph «Flow» (кнопка) | Поток | анимация по рёбрам |
| Due for review (узел графа) | Пора повторить | совпадает с «Reviews due» |
| Related (карточка глоссария) | Связанное | |
| Contact admins | Связаться с администраторами; в узком меню — «Связь с админами» | |
| Beta | бета | |

### Уровни (серверное поле `level`)

| EN | RU |
|---|---|
| Novice | Новичок |
| Apprentice | Подмастерье |
| Student | Ученик |
| Scholar | Учёный |
| Adept | Адепт |
| Expert | Эксперт |
| Master | Мастер |
| Grandmaster | Гроссмейстер |
| Math Architect | Архитектор математики |
| LV | Ур. |

### Технические термины

| EN | RU |
|---|---|
| computerised adaptive testing (CAT) | компьютерное адаптивное тестирование (CAT) |
| covering set | покрывающее множество |
| grader router | маршрутизатор проверки |
| information gain | прирост информации |
| first pass / second pass | с первой попытки / со второй попытки |

## 9. Решения после рецензии первого раунда (2026-10-23)

Приняты пользователем по отчёту `docs/i18n/reviews/2026-10-23-ru-full.md`. Текст уже приведён к ним.

| EN | RU | Комментарий |
|---|---|---|
| prerequisite | необходимая тема; заголовки — «Что изучить сначала» | см. §3 |
| whole number (в дробях) | натуральное число | в арифметике — «целое неотрицательное число», см. §6 |
| XP | множественное число: «XP начисляются», «ваши XP» | |
| starting point | отправная точка; в заголовках и кнопках — «с чего начать» | |
| rate (в задачах на отношения) | отношение величин; unit rate — величина на единицу | «скорость» — о движении; в задачах моделирования (rate × time) — «скорость изменения» |
| isolate (переменную, корень, модуль) | выразить переменную / оставить в одной части; уединить корень или модуль | не «выделите переменную» |
| codomain / range | область прибытия / множество значений | уроки различают эти понятия |
| reference angle | опорный угол | «опорный» в интерфейсе больше не означает prerequisite, поэтому конфликта нет |
| sign analysis | метод интервалов — и в шагах уроков тоже | см. §7 |
| reference points (в условиях) | контрольные точки | не «опорные точки» |

## 10. Решения после рецензии вузовских шагов уроков (2026-10-23)

Приняты пользователем по отчёту `docs/i18n/reviews/2026-10-23-ru-university-steps.md`.

| EN | RU | Комментарий |
|---|---|---|
| modular arithmetic | модульная арифметика | не «модулярная»; так же в украинском |
| FTC1 / FTC2 | первая / вторая основная теорема математического анализа | не «теорема о производной интеграла по верхнему пределу» |
| Worked … (подпись шага) | «Разбор: …»; Worked example(s) — «Разбор примера / примеров» | не «Разобранное сравнение» |
| mini-batch | мини-батч | не «мини-пакет» |

## 11. Решения после рецензии пилота текстов уроков (2026-10-23)

Приняты пользователем по отчёту `docs/i18n/reviews/2026-10-23-ru-lessons-pilot.md`.

| EN | RU | Комментарий |
|---|---|---|
| complementary / supplementary angles | углы, дополняющие друг друга до прямого / до развёрнутого | |
| linear pair | смежные углы | |
| rate (в моделировании: rate × time) | скорость изменения | |
| set-builder notation / interval notation | запись множества условием / запись промежутков | |
| rise / run (наклон прямой) | изменение по вертикали / изменение по горизонтали | решено 2026-10-24: так и в `\text{}` (там нельзя вставить `$x$` и `$y$`), и в прозе уроков наклона; два варианта в одном уроке хуже одного. «Change in $y$ / $x$» из самого оригинала переводится буквально («изменение по $y$ / по $x$») |
| rise / run (лестница) | высота / глубина ступени | |
| in this book | в этом курсе | уроки — часть курса, а не книги |
| larger (о числах со знаком в задачах на сложение) | больше по величине | не «по модулю»: модуль вводится позже |
| kite | дельтоид | |
| perfect square | точный квадрат | |
| reciprocals / opposites | взаимно обратные числа / противоположные числа | |
| PEMDAS и другие английские мнемоники | мнемоники, которые есть в оригинале («Please Excuse My Dear Aunt Sally»), остаются английскими в кавычках; название PEMDAS в оригинале выпало, пропуск остаётся («рассмотрите аббревиатуру .»). Мнемоника C/S — факт об английском: «в английском языке **c**omplementary стоит по алфавиту перед **s**upplementary» | уточнено 2026-10-24: переводчик не дописывает то, чего нет в оригинале |
| isosceles («exactly two equal sides» в оригинале) | переводить как в оригинале | мы только переводим (CLAUDE.md) |

## 12. Тексты уроков по дробям (2026-10-23)

Термины, выбранные при переводе уроков раздела «Дроби» (OpenStax Prealgebra, главы 4 и 5.3, и свои уроки проекта).

| EN | RU | Комментарий |
|---|---|---|
| fraction circles / fraction tiles | круги для дробей / полоски для дробей | модели из набора наглядных пособий |
| manipulatives; Manipulative Mathematics | наглядные пособия; «Математика с наглядными пособиями» | заголовок вставок OpenStax |
| Learning Objectives / Example / Try It / Solution / How To | Цели обучения / Пример / Попробуйте сами / Решение / Как сделать | подписи OpenStax; в заголовках «Как сделать» вторая строка — команда на «вы» («Упростите дробь.») |
| simplified fraction; in simplified form; the fraction is simplified | несократимая дробь; в несократимом виде; дробь несократима | решено 2026-10-24: «упрощённая дробь» не термин, а определение несократимой дроби. «Simplify: 10/15» — «Сократите:», шаг «Simplify the fraction» — «Сократите дробь». «Упростите» остаётся там, где упрощается выражение: числитель, знаменатель, результат умножения, многоэтажная дробь. Дробь с переменной тоже «несократимая» ($\frac{y}{3}$) |
| LCD / LCM | НОЗ (наименьший общий знаменатель) / НОК | аббревиатуры вводятся в скобках при первом упоминании в уроке. «НОЗ» в российских учебниках редкость (пишут полностью), но при 90 упоминаниях на урок допустима |
| convert (rewrite) fractions to equivalent fractions with the LCD | привести дроби к наименьшему общему знаменателю (к НОЗ) | не «привести к равным дробям с НОЗ» и не «перепишите в виде равных дробей» — это калька; рецензия 2026-10-23 |
| whole number (целая часть смешанного числа: «multiply the whole number by the denominator», «add the whole numbers») | целая часть | «натуральное число» — только о числе самом по себе; при сложении и вычитании смешанных чисел складывают «целые части» и «дробные части» |
| Draw a picture / Show a diagram to illustrate your model | Поясните свою модель рисунком | не «нарисуйте/изобразите рисунок» |
| sharing (equally) | деление поровну | не «раздача» |
| repeating decimal | периодическая десятичная дробь | |
| fraction bar | дробная черта | |
| grouping symbols | знаки группировки | |
| quarter (монета в 25 центов) / dime (10 центов) | квотер / дайм | переводим слова как есть, без пояснений от себя: в 4.2 и 4.5 оригинал сам говорит, сколько центов в монете. Игра слов quarter = ¼ = монета передаётся через «четверть»: «дробь ¼ мы называем четвертью, а четверть доллара — это монета квотер», а не «читаем ¼ как квотер» |
| cross-multiplication (в названии темы) | перекрёстное умножение | как в названии темы; «произведения крест-накрест» остаётся в §7 для самих произведений |
| the fraction «of» a number | дробь от числа; слово «от» = умножение | |
| «the quotient of a and b» | частное чисел a и b; «частное от деления a на b» | в задачах на запись выражений |
| borrow (в смешанных числах) | занять единицу | не «заимствовать» |
| «Please excuse my dear Aunt Sally» | оставлена по-английски в кавычках | перевод только то, что сказано в оригинале |

## 13. Тексты уроков: десятичные дроби, проценты, мелкие школьные разделы (2026-10-24)

Термины, выбранные при переводе OpenStax Prealgebra 5.1–5.3 и 6.1–6.3 (`dec.*`, `pct.*`) и коротких уроков проекта (`arith.*`, `pct.*`, `ratio.*`, `prealg.*`).

| EN | RU | Комментарий |
|---|---|---|
| decimal point (в текстах уроков о десятичных дробях) | десятичная точка | решено 2026-10-24: «поставьте десятичную точку в ответе под десятичными точками в заданных числах», пример при этом `18.36`. Так в `dec.*`, `pct.basics`, `pct.convert.*`, `prealg.measure.metric`, `pct.ops.commission`. В названиях тем и шагах `kp/ru` пока «десятичная запятая» (§7); запятая как знак разряда тысяч (`1,000`) остаётся запятой |
| decimal places | знаки после десятичной точки | |
| place holder (нуль-заполнитель) | «дополнить разряды нулями», «дописать нули» | не «держатель места» |
| whole number (в уроках о десятичных дробях) | целое число | так в шагах `kp/ru`; в уроках арифметики — «целое неотрицательное число», в дробях — «натуральное» (§6, §9) |
| counting numbers | натуральные числа | |
| decimal notation | десятичная запись | |
| «and» при чтении десятичной дроби | «целых»: «пятнадцать целых шестьдесят восемь сотых» | в английском «and» обозначает десятичную точку; по-русски эту роль играет «целых», поэтому фраза «Слово «целых» обозначает десятичную точку» переводит тот же приём |
| «th» в названии разряда («thousandth») | слова «десятая, сотая, тысячная» | окончание «th» по-русски не работает, смысл («это дробь, а не число больше единицы») передан через названия |
| equivalent decimals | равные десятичные дроби | как «равные дроби» (§6) |
| penny / nickel | пенни / никель | как quarter и dime в §12: слова переводятся, пояснений от себя нет |
| percent (определение) | процент — отношение со знаменателем 100; «на сто» | «per cent» оставлено латиницей в кавычках |
| percent increase / decrease | процентное увеличение / процентное уменьшение | |
| base (в задачах на проценты) | база | вместе с *величина* для amount; original amount — исходная величина |
| tip | чаевые | |
| sales tax / tax rate / purchase price | налог с продаж / ставка налога / цена покупки | |
| commission / rate of commission | комиссионные / ставка комиссионных | |
| discount / discount rate / amount of discount / sale price / original price | скидка / ставка скидки / размер скидки / цена со скидкой / первоначальная цена | |
| cost (в `pct.ops.markup`) | закупочная цена | решено 2026-10-24: «себестоимость» — затраты производителя, а магазин покупает товар по закупочной цене; selling price там остаётся «ценой продажи» |
| mark-up / mark-up rate / wholesale price / list price | наценка / ставка наценки / оптовая цена / розничная цена | решено 2026-10-24: «оптовая / розничная цена» — естественная пара. «Цена продажи» остаётся для selling price в `pct.ops.markup` |
| unit price / better buy | цена за единицу / выгодная покупка | |
| conversion factor | переводной коэффициент | не «множитель перевода»; рецензия §12 |
| customary units | единицы английской системы мер | |
| Links to Literacy | «Связь с чтением» | вставка OpenStax о детских книгах |
| \text{…} внутри формул | переводить слова: `\text{или}`, `\text{и}`, `\text{4-летние}`, `\text{100 пенни}` | `\text{C}` и `\text{F}` (шкалы Цельсия и Фаренгейта) остаются латиницей, валидатор их принимает |
| кавычки вокруг формулы, оканчивающейся на ” | “…” (английская пара) только в этом месте | например `$0.8\text{?”}$`: формулу менять нельзя, поэтому открывающая кавычка тоже «“» |
| $…$ вокруг английского текста (дефект исходника) | сохранено побайтно | валидатор сверяет формулы символ в символ; затронуты `data-summary` и одна ячейка в `pct.ops.find_rate` и `pct.ops.discount` |

## 14. Тексты уроков по алгебре, вторая половина (2026-10-24)

Термины, выбранные при переводе `teaching/alg.*`, `authored/alg.*`, `algebrica/powers-radicals-logarithms` и `algebrica/sequences` (OpenStax-подобный ORCCA и страницы Algebrica).

| EN | RU | Комментарий |
|---|---|---|
| factor theorem / remainder theorem | теорема о множителе / теорема об остатке (теорема Безу) | rational root theorem — теорема о рациональных корнях |
| prime polynomial | неприводимый многочлен | |
| perfect square (число, выражение) / perfect square trinomial | точный квадрат / трёхчлен — полный квадрат | заголовок урока «Трёхчлены — полные квадраты» |
| generic rectangle | обобщённый прямоугольник | |
| slope triangle | треугольник наклона | |
| rise / run в «rise over run» | «rise over run» оставлено по-английски как мнемоника; сами слова rise / run — «изменение по вертикали» / «изменение по горизонтали» | в формулах нельзя добавлять `$x$` и `$y$`, поэтому слова, а не буквы; в прозе то же (§11, решено 2026-10-24) |
| FOIL (First, Outer, Inner, Last) | в оригинале название пропало; метки F, O, I, L над скобками в формулах переданы словами `\text{первые}`, `\text{внешние}`, `\text{внутренние}`, `\text{последние}` | в прозе «First, Outer, Inner, Last» с русской расшифровкой, как английская мнемоника (§11). Латинские F, O, I, L в `mathOperatorNames` валидатора не добавляем (решено 2026-10-24): единственное место с такими метками переведено словами, а однобуквенные исключения (`E`, `P`) и так спорные |
| point-slope form | уравнение прямой по точке и угловому коэффициенту | slope-intercept form — уравнение прямой с угловым коэффициентом («вид с угловым коэффициентом»); standard form — общее уравнение, «общий вид» |
| vertical / horizontal line test | признак вертикальной / горизонтальной прямой | решено 2026-10-24 (было «проверка вертикальной прямой», в Algebrica — «тест»): «проверка» читалась как проверка самой прямой. Единообразно в уроках, Algebrica и шагах `kp/ru` |
| compound inequality / three-part inequality | составное неравенство / двойное неравенство | |
| one-to-one function | взаимно однозначная функция | |
| principal / compound interest | первоначальная сумма / сложные проценты | |
| augmented matrix / row-echelon form / reduced row-echelon form | расширенная матрица / ступенчатый вид / улучшенный ступенчатый вид | |
| principal (real) n-th root | арифметический корень n-й степени | |
| common / natural / Napierian logarithm | десятичный / натуральный / неперов логарифм | в `\text{}` десятичный логарифм — `lg`, кологарифм — `\text{colog}`: это обозначение операции, а не проза, по-русски его тоже пишут латиницей; валидатор принимает `colog` |
| AM-GM inequality | неравенство между средним арифметическим и средним геометрическим (AM-GM) | |
| monotone convergence theorem | теорема о сходимости монотонной последовательности | |
| eventually increasing / decreasing | возрастающая / убывающая, начиная с некоторого номера | |
| pointwise / uniform convergence | поточечная / равномерная сходимость | |
| polynomial equation / polynomial function | алгебраическое уравнение / полиномиальная функция | решено 2026-10-24, см. §15; «многочленное уравнение / функция» не говорим. Прилагательное «многочленный» остаётся у выражений («многочленное выражение») |
| единицы в `\text{}`: ft, in, yr, wk, gal, h, mile per hour | фт, дюйм, год, нед, гал, ч, миль в час | без пересчёта; склонение по числу |
| названия работ в «Избранных источниках» | английское название в ссылке и русский перевод в скобках после неё: `[Mathematical Induction](url) (Математическая индукция)`; названия университетов по-русски | решено 2026-10-24, см. §15: по-русски такие названия обычно не переводят, а валидатор уроков не принимает абзац со ссылкой, где нет ни одной кириллической буквы |

## 15. Решения пользователя по отчётам рецензентов (2026-10-24)

Пользователь принял все предложения рецензентов из отчётов `2026-10-23-ru-*` и `2026-10-24-ru-*`. Текст приведён к ним. Строки в §7, §11–§14, которые им противоречили, исправлены на месте.

| EN | RU | Комментарий |
|---|---|---|
| polynomial equation / polynomial function | алгебраическое уравнение / полиномиальная функция | не «многочленное уравнение» и не «многочленная функция»; заголовок «Алгебраические уравнения» |
| zero-product property | свойство нулевого произведения | не «правило нулевого множителя» |
| domain of existence / domain of valid solutions | область допустимых значений (ОДЗ) | не «область существования», не «область допустимых решений» |
| repeated root / solution | кратный (двойной) корень | не «повторяющийся» |
| piecewise defined | задана кусочно | в названии типа — «кусочно заданная функция» (§7) |
| conic | коника / коническое сечение | в вузовском тексте — «коника» |
| обозначения внутри `\text{}` | tg, ctg, arctg, arcctg, sgn, НОК, НОД, colog | русские обозначения допустимы, валидатор их принимает |
| (a), (b), (c), (d) — метки пунктов в прозе | а), б), в), г) | кириллица, без внешней скобки, как в российских учебниках; после метки пробел. В формулах (`f(a)`, `\text{…}`, `$…$`) метки не трогаем. Решено 2026-10-24: 440 меток в 11 файлах (`dec.*`, `pct.basics` и копии, `pct.ops.discount`, `pct.ops.tax_tip`) |
| keep / hold (money, time, sets apart) | выражайте деньги в одних и тех же единицах; не смешивайте проценты с основной суммой; следите, чтобы… | не «держите»: калька. Остальные «Держите…» в шагах `kp/ru` (около 25) пока не тронуты |
| order (decimals and fractions, 5.3) | упорядочивание; «упорядочивать» в целях обучения | не «сравнение»: так называется шаг `dec.convert.to_frac`, и текст `frac.ops.to_decimal` совпадает с `dec.convert.*` |
| distribute (глагол) | умножить на каждое слагаемое; раскрыть скобки | решено 2026-10-24: «распределить 3 на $(x+2)$» по-русски не говорят. «Распределительное свойство» — термин, остаётся; «распределить» только в кавычках, когда текст объясняет название свойства («способ „распределить“ число $a$…»). Distribution (существительное) — «раскрытие скобок» |
| The Graph (служебная метка страниц Algebrica) | Граф | решено 2026-10-24: это карта тем (§3), а не график функции. Остальные служебные метки («Тема», «Средняя», «Требует», «Открывает») переведены как были |
| выпавшие ссылки и имена в ORCCA («In  Section ,», «called .», «the  was about») | пропуск переносится как есть, двойными пробелами: «В  разделе  мы…», «называемая .», «2016 года  составлял около» | решено 2026-10-24 для русского пилота (41 урок), как в украинском: правило «мы только переводим» (CLAUDE.md). Описательные ссылки («в предыдущем разделе»), восстановленные имена (Проксима Центавра, постоянная Планка, федеральный долг США, PEMDAS, названия правил степеней) и исправления ошибок оригинала (радиус квадрата, «модули» под заголовком о корнях, «исключая оба значения») возвращены к оригиналу. Заголовки «Смена знака», «Подход «основание на высоту»», «Модуль и неявная группировка» — как в оригинале |
| decimal point | десятичная точка в текстах уроков о десятичных дробях | см. §13; копии одного текста остаются побайтно одинаковыми |
| название работы в «Избранных источниках» | `[English title](url) (русский перевод)` | 26 строк в 7 страницах Algebrica; университеты и авторы — как были |
