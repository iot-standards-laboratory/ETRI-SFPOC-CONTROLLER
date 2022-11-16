import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import 'package:get/get.dart';

import 'app/routes/app_pages.dart';
import 'constants.dart';

void main() {
  serverAddr = kIsWeb ? '${Uri.base.host}:${Uri.base.port}' : 'localhost:4000';
  // serverAddr = 'localhost:4000';

  runApp(
    GetMaterialApp(
      title: "Application",
      initialRoute: AppPages.INITIAL,
      getPages: AppPages.routes,
    ),
  );
}
