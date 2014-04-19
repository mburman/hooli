//
//  WHOMessage.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHOMessage.h"

@implementation WHOMessage

- (instancetype)initWithMessage:(NSString *)message Author:(NSString *)author Location:(CLLocation *)location {
    if (self = [super init]) {
        self.message = message;
        self.author = author;
        self.location = location;
    }
    return self;
}

@end
