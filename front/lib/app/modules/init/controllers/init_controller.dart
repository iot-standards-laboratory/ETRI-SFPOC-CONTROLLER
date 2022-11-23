import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:front/constants.dart';
import 'package:get/get.dart';
import 'package:http/http.dart' as http;

class InitController extends GetxController {
  //TODO: Implement InitController
  var formKey = GlobalKey<FormState>();
  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    super.onClose();
  }

  Future<String?> init({
    required String edgeAddress,
    required String agentName,
    required String accessToken,
  }) async {
    var url = Uri.http(
      serverAddr,
      '/api/v2/init',
    );

    print(url.toString());

    var resp = await http.delete(
      url,
      headers: <String, String>{"access_token": accessToken},
      body: jsonEncode({
        'edgeAddress': edgeAddress,
        'name': agentName,
      }),
    );

    return resp.statusCode == 200 ? null : resp.body;
  }

  Future<String?> initUpdate({
    required String edgeAddress,
    required String agentName,
    required String accessToken,
  }) async {
    var url = Uri.http(
      serverAddr,
      '/api/v2/init',
    );

    print(url.toString());

    var resp = await http.post(
      url,
      headers: <String, String>{"access_token": accessToken},
      body: jsonEncode({
        'edgeAddress': edgeAddress,
        'name': agentName,
      }),
    );

    return resp.statusCode == 200 ? null : resp.body;
  }
}
