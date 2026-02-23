
import 'package:flutter_test/flutter_test.dart';
import 'package:tbdex_mobile/main.dart';

void main() {
  testWidgets('App renders home screen', (WidgetTester tester) async {
    await tester.pumpWidget(const TbdexApp());
    expect(find.text('Good morning,'), findsOneWidget);
  });
}
