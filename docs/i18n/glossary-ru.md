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

- **Десятичная запятая и точка.** В прозе (названия тем, шаги уроков) говорим «десятичная запятая». В числах из условий, ответов и примеров остаётся точка (`2.5`), как и в разделе 2.
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
| rise / run (наклон прямой) | изменение по y / изменение по x | |
| rise / run (лестница) | высота / глубина ступени | |
| in this book | в этом курсе | уроки — часть курса, а не книги |
| larger (о числах со знаком в задачах на сложение) | больше по величине | не «по модулю»: модуль вводится позже |
| kite | дельтоид | |
| perfect square | точный квадрат | |
| reciprocals / opposites | взаимно обратные числа / противоположные числа | |
| PEMDAS и другие английские мнемоники | оставлять как пример английского правила («В английском языке…») | |
| isosceles («exactly two equal sides» в оригинале) | переводить как в оригинале | мы только переводим (CLAUDE.md) |

## 12. Тексты уроков по дробям (2026-10-23)

Термины, выбранные при переводе уроков раздела «Дроби» (OpenStax Prealgebra, главы 4 и 5.3, и свои уроки проекта).

| EN | RU | Комментарий |
|---|---|---|
| fraction circles / fraction tiles | круги для дробей / полоски для дробей | модели из набора наглядных пособий |
| manipulatives; Manipulative Mathematics | наглядные пособия; «Математика с наглядными пособиями» | заголовок вставок OpenStax |
| Learning Objectives / Example / Try It / Solution / How To | Цели обучения / Пример / Попробуйте сами / Решение / Как сделать | подписи OpenStax; в заголовках «Как сделать» вторая строка — команда на «вы» («Упростите дробь.») |
| simplified fraction; in simplified form | упрощённая дробь; в упрощённом виде | в ответах с переменными нельзя сказать «несократимый», поэтому здесь «упрощённый»; «несократимая дробь» остаётся для lowest terms |
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
